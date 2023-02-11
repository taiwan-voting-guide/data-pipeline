package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

type Legislator struct {
	Addr        string `json:"addr"`
	AreaName    string `json:"areaName"`
	Committee   string `json:"committee"`
	Degree      string `json:"degree"`
	Ename       string `json:"ename"`
	Experience  string `json:"experience"`
	Fax         string `json:"fax"`
	LeaveDate   string `json:"leaveDate"`
	LeaveFlag   string `json:"leaveFlag"`
	LeaveReason string `json:"leaveReason"`
	Name        string `json:"name"`
	OnboardDate string `json:"onboardDate"`
	Party       string `json:"party"`
	PartyGroup  string `json:"partyGroup"`
	PicUrl      string `json:"picUrl"`
	Sex         string `json:"sex"`
	Tel         string `json:"tel"`
	Term        string `json:"term"`
}

type LegislatorPayload struct {
	DataList []Legislator `json:"dataList"`
}

func parseConfig() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("loading .env file failed")
	}
}

func get(url string) []byte {
	log.Printf("Fetching %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func parseLegislator(body []byte) []Legislator {
	var result LegislatorPayload
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	log.Printf("Get %d records", len(result.DataList))
	return result.DataList
}

func getCurrentLegislatorInfo() []Legislator {
	url := "https://data.ly.gov.tw/odw/ID9Action.action?fileType=json"
	return parseLegislator(get(url))
}

func getHistoryLegislatorInfo() []Legislator {
	url := "https://data.ly.gov.tw/odw/ID16Action.action?fileType=json"
	return parseLegislator(get(url))
}

func toFile(filePath string, records []Legislator) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	for _, elem := range records {
		data, _ := json.Marshal(elem)
		if _, err := f.Write(data); err != nil {
			log.Fatal(err)
		}
		f.WriteString("\n")
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	parseConfig()
	records := append(getCurrentLegislatorInfo(), getHistoryLegislatorInfo()...)
	toFile("data/legislators.jsonl", records)
}
