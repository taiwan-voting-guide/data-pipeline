package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// 讀取立委資料夾列表
var paths = [...]string{
	"3屆立委",
	"4屆立委",
	"5屆立委",
	"5屆立委",
	"2004第6屆立法委員",
	"2008立委",
	"20120114-總統及立委",
	"2016總統立委",
	"2020總統立委",
}

// 候選人欄位對應
var candColumns = [...]string{
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

func readCSV(filepath string) ([][]string, error) {
	// Open the CSV file.
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the CSV data.
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func renameFile(filePath string, newFileName string) {
	filePath = fmt.Sprintf("cmd/data_gov_tw/votedata/voteData/%s", filePath)
	newFileName = fmt.Sprintf("cmd/data_gov_tw/votedata/voteData/%s", newFileName)
	_, err := os.Stat(filePath)
	if err == nil {
		err := os.Rename(filePath, newFileName)
		if err != nil {
			fmt.Println("Failed to rename file:", err)
		} else {
			fmt.Println("File renamed successfully.")
		}
	}
}

func main() {
	// 修正檔名
	renameFile("2016總統立委/山地立委/elcand_T3.csv", "2016總統立委/山地立委/elcand.csv")
	renameFile("2016總統立委/平地立委/elcand_T2.csv", "2016總統立委/平地立委/elcand.csv")
	renameFile("2016總統立委/區域立委/elcand_T1.csv", "2016總統立委/區域立委/elcand.csv")

	for _, path := range paths {
		log.Printf("Load /%s ...", path)

		firstFolderPath := fmt.Sprintf("cmd/data_gov_tw/votedata/voteData/%s", path)

		firstFolder, err := os.Open(firstFolderPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer firstFolder.Close()

		secondFolders, err := firstFolder.Readdir(-1)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, secondFolder := range secondFolders {
			folderName := secondFolder.Name()
			if !strings.Contains(folderName, "立委") && !strings.Contains(folderName, "山原") && !strings.Contains(folderName, "平原") && !strings.Contains(folderName, "區域") {
				continue
			}

			log.Printf("/%s", folderName)

			candPath := fmt.Sprintf("%s/%s/elcand.csv", firstFolderPath, folderName)

			// read csv
			records, err := readCSV(candPath)
			if err != nil {
				fmt.Println(err)
				return
			}

			// 創建一個檔案，寫入 JSON 格式的資料
			fileName := fmt.Sprintf("data/%s_%s.jsonl", path, folderName)
			file, err := os.Create(fileName)
			if err != nil {
				panic(err)
			}
			for _, record := range records {
				// 轉成 map 變數
				recordMap := make(map[string]string)
				for i, column := range record {
					recordMap[candColumns[i]] = column
				}

				// 將 map 變數轉換成 JSON 格式
				jsonData, err := json.Marshal(recordMap)
				if err != nil {
					log.Fatal(err)
				}

				// 寫入 JSON 格式的資料
				if _, err := file.Write(jsonData); err != nil {
					log.Fatal(err)
				}
				file.WriteString("\n")
			}
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}
}
