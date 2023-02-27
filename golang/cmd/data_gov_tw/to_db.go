package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/taiwan-voting-guide/backend/candidate"
	"github.com/taiwan-voting-guide/backend/model"
	"github.com/taiwan-voting-guide/backend/politician"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("loading .env file failed")
	}

	importToDB("../data/legislators_history_7.jsonl", 7)
}

func importToDB(path string, term int) {
	datas, _ := readJSONLFile("../data/legislators_history_7.jsonl")

	for _, data := range datas {
		fmt.Println(data)

		// Create Politician
		formatBirthdate := strings.Replace(data["birthday"].(string), "00-00", "01-01", 1)
		politicians, err := politician.New().SearchByNameAndBirthdate(context.Background(), data["name"].(string), formatBirthdate)
		if err != nil {
			fmt.Println(err)
			return
		}

		var politicianId int64
		if len(politicians) != 0 {
			politicianId = (*politicians[0]).Id
		} else {
			politicianId = createPolitician(data)
		}

		// Create Candidate by id
		createCandidate(term, int(politicianId), data)
	}
}

func createCandidate(term int, politicianId int, data map[string]interface{}) int64 {
	elected := data["is_victor"] == "*"
	no, _ := strconv.Atoi(data["no"].(string))
	partyId, _ := strconv.Atoi(data["party_code"].(string))
	c := &model.CandidateLyRepr{
		Type:         getCandidateType(data["base_name"].(string)),
		Term:         term,
		PoliticianId: int(politicianId),
		Number:       no,
		Elected:      elected,
		PartyId:      partyId,
		Area:         data["base_name"].(string),
	}

	err := candidate.New().Create(context.Background(), c.Model())
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return 1
}

func createPolitician(data map[string]interface{}) int64 {
	name := data["name"].(string)
	var birthdate *string
	if data["birthday"] != nil {
		formatBirthdate := strings.Replace(data["birthday"].(string), "00-00", "01-01", 1)
		birthdate = &formatBirthdate
	}
	sex := getPoliticianSex(data["sex"].(string))
	// partyId, _ := strconv.Atoi(data["party_code"].(string))

	p := &model.PoliticianRepr{
		Name:      name,
		Birthdate: birthdate,
		Sex:       sex,
		// PartyId:   partyId, // TODO
	}

	id, err := politician.New().Create(context.Background(), p)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	log.Printf("Create politician: %s(%d)", name)

	return id
}

func getCandidateType(baseName string) model.CandidateType {
	switch baseName {
	case "山地原住民選區":
		return model.CandidateType("ly-mountain")
	case "平地原住民選區":
		return model.CandidateType("ly-ground")
	default:
		return model.CandidateType(model.CandidateTypeLyLocal)
	}
}

func getPoliticianSex(sexCode string) model.Sex {
	switch sexCode {
	case "1":
		return model.Sex(model.SexMale)
	case "2":
		return model.Sex(model.SexFemale)
	}
	return ""
}

func readJSONLFile(filename string) ([]map[string]interface{}, error) {
	// 打開 JSONL 檔案
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 創建 Scanner 對象
	scanner := bufio.NewScanner(file)

	// 逐行讀取 JSONL 檔案
	var result []map[string]interface{}
	for scanner.Scan() {
		// 解碼 JSON
		var data map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &data); err != nil {
			return nil, err
		}

		// 將解碼後的 JSON 對象添加到結果中
		result = append(result, data)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
