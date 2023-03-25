package util

import (
	"encoding/csv"
	"os"
)

func ReadCSV(filepath string) ([][]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func ReadCSVToMap(filepath string, columnsMapper []string) ([]map[string]string, error) {
	records, err := ReadCSV(filepath)
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

func CreateStagingEndpoint() string {
	backendEndpoint := os.Getenv("BACKEND_URL")
	return backendEndpoint + "/workspace/staging/create"
}
