// 更新politicians資料表
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/taiwan-voting-guide/backend/model"

	"github.com/taiwan-voting-guide/data/cec"
	"github.com/taiwan-voting-guide/data/util"
)

func generateStagingCreates() ([]model.Staging, error) {
	stagingCreate := []model.Staging{}
	for _, fileNames := range cec.LyCandFileNames {
		cands, _, err := fileNames.Load()
		if err != nil {
			return nil, err
		}

		for _, cand := range cands {
			staging := model.Staging{
				Table: model.StagingTablePoliticians,
				SearchBy: model.StagingFields{
					"name":      cand.Name,
					"birthdate": cand.Birthdate,
				},
				Fields: model.StagingFields{
					"name":      cand.Name,
					"birthdate": cand.Birthdate,
					"sex":       cand.Sex,
				},
			}

			if *cand.PartyId != 999 {
				staging.Fields["current_party_id"] = *cand.PartyId
			}

			stagingCreate = append(stagingCreate, staging)
		}
	}

	return stagingCreate, nil
}

func main() {
	godotenv.Load()
	stagings, _ := generateStagingCreates()
	for _, staging := range stagings {
		stagingCreateJson, err := json.Marshal(staging)
		if err != nil {
			log.Fatal()
		}

		resp, err := http.Post(util.CreateStagingEndpoint(), "application/json", bytes.NewReader(stagingCreateJson))
		if err != nil {
			log.Println(err)
		}

		if resp.StatusCode != 201 {
			log.Printf("staging create failed: %s", stagingCreateJson)
		}
	}
}
