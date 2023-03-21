// https://data.cec.gov.tw/選舉資料庫/votedata.zip
package cec

import (
	"fmt"
	"path"
	"strconv"

	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/data/util"
)

const rootPath = "../data/votedata/voteData"

type LyCandFileName struct {
	Cand string
	Base string
}

func (fn LyCandFileName) Load() ([]Cand, []Base, error) {
	maps, err := util.ReadCSVToMap(fn.Cand, CandColumns)
	if err != nil {
		return nil, nil, err
	}

	cands := []Cand{}
	for _, m := range maps {
		cands = append(cands, Cand{
			PrvCode:    m["prv_code"],
			CityCode:   m["city_code"],
			AreaCode:   m["area_code"],
			DeptCode:   m["dept_code"],
			LiCode:     m["li_code"],
			No:         m["no"],
			Name:       m["name"],
			PartyId:    partyId(m["party_code"]),
			Sex:        sex(m["sex"]),
			Birthdate:  birthdate(m["birthday"]),
			Age:        m["age"],
			Birthplace: m["birthplace"],
			Degree:     m["degree"],
			IsCurrent:  m["is_current"],
			IsVictor:   m["is_victor"],
			Vice:       m["vice"],
		})
	}

	maps, err = util.ReadCSVToMap(fn.Base, BaseColumns)
	if err != nil {
		return nil, nil, err
	}

	bases := []Base{}
	for _, m := range maps {
		bases = append(bases, Base{
			PrvCode:  m["prv_code"],
			CityCode: m["city_code"],
			AreaCode: m["area_code"],
			DeptCode: m["dept_code"],
			LiCode:   m["li_code"],
			Name:     m["name"],
		})
	}

	return cands, bases, nil
}

// 立委檔案列表
var LyCandFileNames = [...]LyCandFileName{
	// 第10屆立委
	{
		path.Join(rootPath, "2020總統立委/平地立委/elcand.csv"),
		path.Join(rootPath, "2020總統立委/平地立委/elbase.csv"),
	},
	{
		path.Join(rootPath, "2020總統立委/山地立委/elcand.csv"),
		path.Join(rootPath, "2020總統立委/山地立委/elbase.csv"),
	},
	{
		path.Join(rootPath, "2020總統立委/區域立委/elcand.csv"),
		path.Join(rootPath, "2020總統立委/區域立委/elbase.csv"),
	},
	// 第9屆立委
	{
		path.Join(rootPath, "2016總統立委/山地立委/elcand_T3.csv"),
		path.Join(rootPath, "2016總統立委/山地立委/elbase_T3.csv"),
	},
	{
		path.Join(rootPath, "2016總統立委/平地立委/elcand_T2.csv"),
		path.Join(rootPath, "2016總統立委/平地立委/elbase_T2.csv"),
	},
	{
		path.Join(rootPath, "2016總統立委/區域立委/elcand_T1.csv"),
		path.Join(rootPath, "2016總統立委/區域立委/elbase_T1.csv"),
	},
	// TODO
	// 20120114-總統及立委
	// 2008立委
	// 2004第6屆立法委員
	// 5屆立委
	// 4屆立委
	// 3屆立委
}

// 候選人欄位對應
type Cand struct {
	PrvCode    string
	CityCode   string
	AreaCode   string
	DeptCode   string
	LiCode     string
	No         string
	Name       string
	PartyId    *int
	Sex        model.Sex
	Birthdate  string
	Age        string
	Birthplace string
	Degree     string
	IsCurrent  string
	IsVictor   string
	Vice       string
}

var CandColumns = []string{
	"prv_code",   // 省市別
	"city_code",  // 縣市別
	"area_code",  // 選區別
	"dept_code",  // 鄉鎮市區
	"li_code",    // 村里別
	"no",         // 號次
	"name",       // 名字
	"party_code", // 政黨代號
	"sex",        // 性別
	"birthday",   // 出生日期
	"age",        // 年齡
	"birthplace", // 出生地
	"degree",     // 學歷
	"is_current", // 現任
	"is_victor",  // 當選註記
	"vice",       // 副手
}

// 行政區欄位對應
type Base struct {
	PrvCode  string
	CityCode string
	AreaCode string
	DeptCode string
	LiCode   string
	Name     string
}

var BaseColumns = []string{
	"prv_code",  // 省市別
	"city_code", // 縣市別
	"area_code", // 選區別
	"dept_code", // 鄉鎮市區
	"li_code",   // 村里別
	"name",      // 行政區名稱
	"pass",      // 跳過
	"pass",      // 跳過
}

func sex(sexCode string) model.Sex {
	switch sexCode {
	case "1":
		return model.Sex(model.SexMale)
	case "2":
		return model.Sex(model.SexFemale)
	}
	return ""
}

func birthdate(str string) string {
	year := 1911 + int(str[0]-'0')*100 + int(str[1]-'0')*10 + int(str[2]-'0')
	month := int(str[3]-'0')*10 + int(str[4]-'0')
	day := int(str[5]-'0')*10 + int(str[6]-'0')

	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func partyId(partyCode string) *int {
	i, err := strconv.Atoi(partyCode)
	if err != nil {
		return nil
	}

	return &i
}
