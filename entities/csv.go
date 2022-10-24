package entities

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sync"
)

type CSVBlock struct {
	Header []string
	Data   []map[string]interface{}
}

type CSVMap map[string]*CSVBlock

type CSVMapWithMutex struct {
	sync.RWMutex
	Map CSVMap
}

type LabelIndex struct {
	Label string
	Index int
}

func contains(ss []string, s string) bool {
	for _, d := range ss {
		if d == s {
			return true
		}
	}
	return false
}

func (cwm CSVMapWithMutex) Add(filename string, csm *CSVBlock) {
	cwm.Lock()
	defer cwm.Unlock()
	cwm.Map[filename] = csm
}

func FilterField(data []map[string]interface{}, allowed []string) (result []map[string]interface{}) {
	result = []map[string]interface{}{}
	for _, dataCSV := range data {
		tempResult := map[string]interface{}{}
		for _, allowedLabel := range allowed {
			d := dataCSV[allowedLabel]
			if d != nil {
				tempResult[allowedLabel] = d
			}
		}
		result = append(result, tempResult)
	}
	return
}

func QueryFilterContent(data []map[string]interface{}, query, label, value string) (result []map[string]interface{}) {
	for _, datum := range data {
		switch query {
		case "!=":
			if datum[label] != value {
				result = append(result, datum)
			}
		case "=":
			if datum[label] == value {
				result = append(result, datum)
			}
		default:
			continue
		}
	}
	return
}

func NewCSVBlockFromFile(filePath string) (*CSVBlock, error) {
	// Open CSV file
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Unable to parse file as CSV for %s : %s", filePath, err.Error()),
		)
	}

	// Parse CSV file stream to string of arrays
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Unable to parse file as CSV for %s : %s", filePath, err.Error()),
		)
	}

	// Convert all CSV string array to array of HashMap to ease to call the value.
	headers := records[0]
	result := []map[string]interface{}{}

	for _, record := range records[1:] {
		tempResult := map[string]interface{}{}
		for i, header := range headers {
			tempResult[header] = record[i]
		}
		result = append(result, tempResult)
	}

	// set-in into data
	data := CSVBlock{
		Header: headers,
		Data:   result,
	}

	return &data, nil
}
