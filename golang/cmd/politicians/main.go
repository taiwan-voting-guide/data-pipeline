// 更新politicians資料表
package main

import (
	"log"

	"github.com/taiwan-voting-guide/backend/model"

	"github.com/taiwan-voting-guide/data/cec"
)

func generateStagingCreates() ([]model.Staging, error) {
	stagingCreate := []model.Staging{}
	for _, fileNames := range cec.LyCandFileNames {
		cands, _, err := fileNames.Load()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		for _, cand := range cands {
			stagingCreate = append(stagingCreate, model.Staging{
				Table: model.StagingTablePoliticians,
				SearchBy: model.StagingFields{
					"name":      cand.Name,
					"birthdate": cand.Birthdate,
				},
				Fields: model.StagingFields{
					"name":             cand.Name,
					"birthdate":        cand.Birthdate,
					"sex":              cand.Sex,
					"current_party_id": cand.PartyId,
				},
			})
		}
	}

	return stagingCreate, nil
}

func main() {
	stagings, _ := generateStagingCreates()
	for i := 0; i < 10; i++ {
		log.Println(stagings[i])
	}
}
