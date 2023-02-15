package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
)

// 資料路徑
var sourcePath = "data/votedata"

// 立委資料夾列表 對應 存擋檔名
var paths = map[string]string{
	"3屆立委": "legislators_history_3",
	"4屆立委": "legislators_history_4",
	"5屆立委": "legislators_history_5",
	"2004第6屆立法委員": "legislators_history_6",
	"2008立委": "legislators_history_7",
	"20120114-總統及立委": "legislators_history_8",
	"2016總統立委": "legislators_history_9",
	"2020總統立委": "legislators_history_10",
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

// 行政區欄位對應
var baseColumns = [...]string{
	"prv_code",   // 省市別
	"city_code",  // 縣市別
	"area_code",  // 選區別
	"dept_code",  // 鄉鎮市區
	"li_code",    // 村里別
	"name",       // 行政區名稱
	"pass",       // 跳過
	"pass",       // 跳過
}

// 政黨欄位對應
var patyColumns = [...]string{
	"party_code", // 政黨代號
	"party",      // 政黨名稱
}

func readCSV(filepath string) ([][]string, error) {
	// 打開 CSV 檔案
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 讀取 CSV 資料
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func readCSVToMap(filepath string, columnsMapper []string) ([]map[string]string, error) {
	records, err := readCSV(filepath)
	if err != nil {
		return nil, err
	}

	var result []map[string]string
	for _, record := range records {
		recordMap := make(map[string]string)
		for i, column := range record {
			recordMap[columnsMapper[i]] = column
		}
		result = append(result, recordMap)
	}

	return result, nil
}

func findSimilarMap(data []map[string]string, input map[string]string) (map[string]string, bool) {
	for _, item := range data {
		match := true
		for key, value := range input {
			if item[key] != value {
				match = false
				break
			}
		}
		if match {
			return item, true
		}
	}
	return nil, false
}


func renameFile(filePath string, newFileName string) {
	filePath = fmt.Sprintf("%s/voteData/%s", sourcePath, filePath)
	newFileName = fmt.Sprintf("%s/voteData/%s", sourcePath, newFileName)
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

func runFirstFolders() {
	for path, filename := range paths {
		log.Printf("Load /%s ...", path)

		firstFolderPath := fmt.Sprintf("%s/voteData/%s", sourcePath, path)

		firstFolder, err := os.Open(firstFolderPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer firstFolder.Close()

		// 取得子資料夾
		secondFolders, err := firstFolder.Readdir(-1)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 創建一個檔案，寫入 JSON 格式的資料
		fileName := fmt.Sprintf("data/%s.jsonl", filename)
		file, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}

		// 讀取子資料夾
		runSecondFolders(firstFolderPath, file, secondFolders)

		// 關閉檔案
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func runSecondFolders(firstFolderPath string, file *os.File, secondFolders []os.FileInfo) {
	for _, secondFolder := range secondFolders {
		folderName := secondFolder.Name()

		// 只需要讀取部分資料夾(山地、平地、區域)
		if !strings.Contains(folderName, "立委") && !strings.Contains(folderName, "山原") && !strings.Contains(folderName, "平原") && !strings.Contains(folderName, "區域") {
			continue
		}

		log.Printf("/%s", folderName)

		// 檔案路徑
		candPath := fmt.Sprintf("%s/%s/elcand.csv", firstFolderPath, folderName)
		basePath := fmt.Sprintf("%s/%s/elbase.csv", firstFolderPath, folderName)
		patyPath := fmt.Sprintf("%s/%s/elpaty.csv", firstFolderPath, folderName)

		// 讀取 csv
		candRecords, err := readCSV(candPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		baseDatas, err := readCSVToMap(basePath, baseColumns[:])
		if err != nil {
			fmt.Println(err)
			return
		}
		patyDatas, err := readCSVToMap(patyPath, patyColumns[:])
		if err != nil {
			fmt.Println(err)
			return
		}

		// 讀取欄位
		for _, record := range candRecords {
			// 轉成 map 變數
			recordMap := make(map[string]interface{})

			// format 欄位
			for i, column := range record {
				recordName := candColumns[i]
				trimColumn := strings.Trim(column, " ")
				trimColumn = strings.Trim(trimColumn, "'")

				if recordName == "age" {
					ageColumn, err := strconv.Atoi(trimColumn)
					if err != nil {
						fmt.Println(err)
					}
					recordMap[recordName] = ageColumn
				} else if recordName == "birthday" {
					yearNumber, _ := strconv.Atoi(trimColumn[:3])
					yearNumber += 1911
					isOnlyYear := len(trimColumn) == 3
					if isOnlyYear {
						recordMap[recordName] = fmt.Sprintf("%s-00-00", strconv.Itoa(yearNumber))
					} else {
						recordMap[recordName] = fmt.Sprintf("%s-%s-%s", strconv.Itoa(yearNumber), trimColumn[3:5], trimColumn[5:7])
					}

				} else {
					recordMap[recordName] = trimColumn
				}
			}

			// 取得行政區欄位
			if strings.Contains(folderName, "山原") || strings.Contains(folderName, "山地") {
				recordMap["base_name"] = "山地原住民選區"
			} else if strings.Contains(folderName, "平原") || strings.Contains(folderName, "平地") {
				recordMap["base_name"] = "平地原住民選區"
			} else {
				findBase := make(map[string]string)
				findBase["prv_code"] = recordMap["prv_code"].(string)
				findBase["city_code"] = recordMap["city_code"].(string)
				findBase["area_code"] = recordMap["area_code"].(string)
				findBase["dept_code"] = recordMap["dept_code"].(string)
				findBase["li_code"] = recordMap["li_code"].(string)
				baseData, isMatch := findSimilarMap(baseDatas, findBase)
				if isMatch {
					recordMap["base_name"] = baseData["name"]
				} else {
					recordMap["base_name"] = ""
				}
			}

			// 取得政黨欄位
			findParty := make(map[string]string)
			findParty["party_code"] = recordMap["party_code"].(string)
			partyData, isMatch := findSimilarMap(patyDatas, findParty)
			if isMatch {
				recordMap["party"] = partyData["party"]
			} else {
				recordMap["party"] = ""
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
	}
}

func fixFilenames() {
	renameFile("2016總統立委/山地立委/elbase_T3.csv", "2016總統立委/山地立委/elbase.csv")
	renameFile("2016總統立委/平地立委/elbase_T2.csv", "2016總統立委/平地立委/elbase.csv")
	renameFile("2016總統立委/區域立委/elbase_T1.csv", "2016總統立委/區域立委/elbase.csv")
	renameFile("2016總統立委/山地立委/elcand_T3.csv", "2016總統立委/山地立委/elcand.csv")
	renameFile("2016總統立委/平地立委/elcand_T2.csv", "2016總統立委/平地立委/elcand.csv")
	renameFile("2016總統立委/區域立委/elcand_T1.csv", "2016總統立委/區域立委/elcand.csv")
}

func main() {
	fixFilenames()
	runFirstFolders()
}
