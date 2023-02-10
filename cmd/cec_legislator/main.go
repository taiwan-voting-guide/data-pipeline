package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// 選舉主題
type Theme struct {
	ThemeID          string `json:"theme_id"`
	ThemeGroup       string `json:"theme_group"`
	TypeID           string `json:"type_id"`
	SubjectID        string `json:"subject_id"`
	LegislatorTypeId string `json:"legislator_type_id"`
	DataLevel        string `json:"data_level"`
	ThemeName        string `json:"theme_name"`
	VoteDate         string `json:"vote_date"`
	LegislatorDesc   string `json:"legislator_desc"`
	HasData          bool   `json:"has_data"`
}

type Area struct {
	AreaName   string  `json:"area_name"`
	ThemeItems []Theme `json:"theme_items"`
}

func fetchURL(url string) []byte {
	log.Printf(fmt.Sprintf("Fetching %s", url))

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching URL: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
	}

	return body
}

/*
只抓取立法委員資料

資料來源：中選會-選舉資料庫
網址：https://db.cec.gov.tw/ElecTable/Election?type=Legislator
*/
func main() {
	// A. load body from url
	// body := fetchURL("https://db.cec.gov.tw/static/elections/list/ELC_L0.json")

	// to json
	// err := ioutil.WriteFile("output.json", body, 0644)
	// if err != nil {
	//     log.Fatalf("Error writing to file: %s", err)
	// }

	// B. load body from json
	body, err := ioutil.ReadFile("output.json")
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	// parse body
	var areas []Area
	err = json.Unmarshal(body, &areas)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}

	// 選區
	for _, area := range areas {
		fmt.Println("Area Name:", area.AreaName)
		for _, theme := range area.ThemeItems {
			var url string
			switch theme.LegislatorTypeId {
			case "L1": // 區域(列表)
				url = fmt.Sprintf("https://db.cec.gov.tw/static/elections/data/areas/%s/%s/%s/%s/C/00_000_00_000_0000.json", theme.TypeID, theme.SubjectID, theme.LegislatorTypeId, theme.ThemeID)
				// TODO: 抓區域連結 > 抓候選人清單

			case "L2": // 平地原住民(候選人)
				url = fmt.Sprintf("https://db.cec.gov.tw/static/elections/data/tickets/%s/%s/%s/%s/N/00_000_00_000_0000.json", theme.TypeID, theme.SubjectID, theme.LegislatorTypeId, theme.ThemeID)
				// TODO: 抓候選人清單

			case "L3": // 山地原住民(候選人)
				url = fmt.Sprintf("https://db.cec.gov.tw/static/elections/data/tickets/%s/%s/%s/%s/N/00_000_00_000_0000.json", theme.TypeID, theme.SubjectID, theme.LegislatorTypeId, theme.ThemeID)
				// TODO: 抓候選人清單

			case "L4": // 不分區政黨(政黨)，跳過
				continue
			}

			title := fmt.Sprintf("%s_%s_%s", theme.ThemeName, theme.VoteDate, theme.LegislatorDesc)
			log.Printf(fmt.Sprintf("%s: %s", title, url))
		}
	}
}
