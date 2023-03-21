package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/extrame/xls"
	"github.com/joho/godotenv"

	"github.com/taiwan-voting-guide/backend/model"
)

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
		chromedp.Navigate("https://party.moi.gov.tw/PartyMain.aspx?n=16100&sms=13073"),
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

	stagingCreates := []model.Staging{}
	if xlFile, err := xls.Open(tmpDir+"/"+filename, "utf-8"); err == nil {
		if sheet := xlFile.GetSheet(0); sheet != nil {
			for row := 1; row <= int(sheet.MaxRow); row++ {
				fmt.Println(row)
				stagingCreates = append(stagingCreates, rowToStagingCreate(sheet.Row(row)))
			}
		}
	}

	stagingCreates = stagingCreates[1:]

	// http request to staging api
	for _, stagingCreate := range stagingCreates {
		stagingCreateJson, err := json.Marshal(stagingCreate)
		if err != nil {
			log.Fatal()
		}

		// get env config
		backendEndpoint := os.Getenv("BACKEND_ENDPOINT")
		endpoint := backendEndpoint + "/workspace/staging"

		log.Println(string(stagingCreateJson))

		http.Post(endpoint, "application/json", bytes.NewReader(stagingCreateJson))
	}
}

func rowToStagingCreate(row *xls.Row) model.Staging {
	fmt.Println("row to stg")
	id, _ := strconv.Atoi(row.Col(0))
	fmt.Println("id")
	chairman := ""
	if !strings.Contains(row.Col(2), "負責人") {
		chairman = row.Col(2)
	}

	return model.Staging{
		Table: "parties",
		SearchBy: model.StagingFields{
			"id": id,
		},
		Fields: model.StagingFields{
			"id":                  id,
			"name":                row.Col(1),
			"chairman":            chairman,
			"established_date":    ROCDateToDate(row.Col(3)),
			"filing_date":         ROCDateToDate(row.Col(4)),
			"main_office_address": row.Col(5),
			"mailing_address":     row.Col(6),
			"phone_number":        row.Col(7),
			"status":              row.Col(8),
		},
	}

}

func ROCDateToDate(date string) string {
	year := 0
	month := 0
	day := 0

	date = strings.TrimSpace(date)
	date = strings.Replace(date, "前", "-", 1)
	_, err := fmt.Sscanf(date, "民國%d年%d月%d日", &year, &month, &day)
	if err != nil {
		log.Println(err)
		return ""
	}

	year += 1911

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02T15:04:05Z")
}
