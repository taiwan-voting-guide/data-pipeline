package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/extrame/xls"
	"github.com/joho/godotenv"
	"github.com/taiwan-voting-guide/backend/pg"
)

type Party struct {
	Id                int       `json:"id"`
	Name              string    `json:"name"`
	Chairman          string    `json:"chairman"`
	EstablishedDate   time.Time `json:"established_date"`
	FilingDate        time.Time `json:"filing_date"`
	MainOfficeAddress string    `json:"main_office_address"`
	MailingAddress    string    `json:"mailing_address"`
	PhoneNumber       string    `json:"phone_number"`
	Status            int       `json:"status"`
}

type Record struct {
	Table  string `json:"table"`
	Record Party  `json:"record"`
}

type Records []Record

func main() {
	godotenv.Load()

	tmpDir, err := ioutil.TempDir("", "chromedp-")
	if err != nil {
		log.Fatal(err)
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var nodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://party.moi.gov.tw"),
		chromedp.Nodes("search_party", &nodes, chromedp.ByID),
	); err != nil {
		log.Println(err)
	}

	done := make(chan string, 1)
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if ev, ok := ev.(*browser.EventDownloadProgress); ok {
			if ev.State == browser.DownloadProgressStateCompleted {
				done <- ev.GUID
				close(done)
			}
		}

	})

	var ns []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.MouseClickNode(nodes[0]),
		chromedp.WaitVisible("ContentPlaceHolder1_BTN_Export_Excel", chromedp.ByID),
		chromedp.Nodes("ContentPlaceHolder1_BTN_Export_Excel", &ns, chromedp.ByID),
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).WithDownloadPath(tmpDir).WithEventsEnabled(true),
	); err != nil {
		log.Println(err)
	}

	filename := ""
	if err := chromedp.Run(ctx,
		chromedp.MouseClickNode(ns[0]),
		chromedp.ActionFunc(func(ctx context.Context) error {
			filename = <-done
			return nil
		}),
	); err != nil {
		log.Println(err)
	}

	recordsList := []Records{}

	if xlFile, err := xls.Open(tmpDir+"/"+filename, "utf-8"); err == nil {
		if sheet := xlFile.GetSheet(0); sheet != nil {
			for row := 0; row <= int(sheet.MaxRow); row++ {
				recordsList = append(recordsList, Records{Record{
					Table:  "parties",
					Record: rowToParty(sheet.Row(row)),
				}})
			}
		}
	}

	conn, err := pg.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	for _, r := range recordsList {
		recordJson, err := json.Marshal(r)
		if err != nil {
			log.Fatal(err)
		}

		_, err = conn.Exec(context.Background(), "INSERT INTO staging_data (records) VALUES ($1)", recordJson)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func rowToParty(row *xls.Row) Party {
	id, _ := strconv.Atoi(row.Col(0))
	chairman := ""
	if !strings.Contains(row.Col(2), "負責人") {
		chairman = row.Col(2)
	}

	return Party{
		Id:                id,
		Name:              row.Col(1),
		Chairman:          chairman,
		EstablishedDate:   ROCDateToDate(row.Col(3)),
		FilingDate:        ROCDateToDate(row.Col(4)),
		MainOfficeAddress: row.Col(5),
		MailingAddress:    row.Col(6),
		PhoneNumber:       row.Col(7),
		Status:            statusStrToNum(row.Col(8)),
	}
}

func statusStrToNum(status string) int {
	switch status {
	case "一般":
		return 1
	case "撤銷備案":
		return 2
	case "自行解散":
		return 3
	case "失聯":
		return 4
	case "廢止備案":
		return 5
	}

	return 0
}

func ROCDateToDate(date string) time.Time {
	year := 0
	month := 0
	day := 0

	date = strings.TrimSpace(date)
	date = strings.Replace(date, "前", "-", 1)
	_, err := fmt.Sscanf(date, "民國%d年%d月%d日", &year, &month, &day)
	if err != nil {
		log.Println(err)
		return time.Time{}
	}

	year += 1911

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
