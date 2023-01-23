package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

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

type CurrentLegislatorPayload struct {
	JsonList []Legislator `json:"jsonList"`
}

type HistoryLegislatorPayload struct {
	DataList []Legislator `json:"dataList"`
}

func parseConfig() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("loading .env file failed")
	}
}

func getLegislatorInfo(url string) []byte {
	log.Println(fmt.Sprintf("Fetching %s", url))
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

func getCurrentLegislatorInfo() chan Legislator {
	url := "https://data.ly.gov.tw/odw/openDatasetJson.action?id=9&selectTerm=all"
	body := getLegislatorInfo(url)
	var result CurrentLegislatorPayload
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	ch := make(chan Legislator)
	go func(c chan Legislator) {
		for _, legislator := range result.JsonList {
			log.Println(legislator)
			c <- legislator
		}
		close(ch)
	}(ch)
	log.Println(fmt.Sprintf("Get %d records from %s", len(result.JsonList), url))
	return ch
}

func getTermLegislatorInfo(term int) []Legislator {
	url := fmt.Sprintf("https://data.ly.gov.tw/odw/ID16Action.action?term=%02d&fileType=json", term)
	body := getLegislatorInfo(url)
	var result HistoryLegislatorPayload
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	log.Println(fmt.Sprintf("Get %d records from %s", len(result.DataList), url))
	return result.DataList
}

func getHistoryLegislatorInfo() chan Legislator {
	log.Println("Get history legislator info")
	ch := make(chan Legislator)
	wg := sync.WaitGroup{}
	for i := 1; i <= viper.GetInt("legislator.last_term"); i++ {
		wg.Add(1)
		go func(term int) {
			defer wg.Done()
			for _, legislator := range getTermLegislatorInfo(term) {
				ch <- legislator
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	log.Println("End history legislator info")
	return ch
}

func main() {
	parseConfig()
	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	for elem := range getCurrentLegislatorInfo() {
		data, _ := json.Marshal(elem)
		if _, err := f.Write(data); err != nil {
			log.Fatal(err)
		}
		f.WriteString("\n")
	}
	for elem := range getHistoryLegislatorInfo() {
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
