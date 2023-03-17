// https://data.cec.gov.tw/選舉資料庫/votedata.zip
package cec

// 立委資料夾列表 對應 存擋檔名
var FolderNames = map[string]string{
	"3屆立委":           "legislators_history_3",
	"4屆立委":           "legislators_history_4",
	"5屆立委":           "legislators_history_5",
	"2004第6屆立法委員":    "legislators_history_6",
	"2008立委":         "legislators_history_7",
	"20120114-總統及立委": "legislators_history_8",
	"2016總統立委":       "legislators_history_9",
	"2020總統立委":       "legislators_history_10",
}

// 候選人欄位對應
var CandColumns = [...]string{
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
var BaseColumns = [...]string{
	"prv_code",  // 省市別
	"city_code", // 縣市別
	"area_code", // 選區別
	"dept_code", // 鄉鎮市區
	"li_code",   // 村里別
	"name",      // 行政區名稱
	"pass",      // 跳過
	"pass",      // 跳過
}
