package dataengine

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type CSVEngine struct {
	RootDir string
}

func (e *CSVEngine) Read(fname string) ([]map[string]any, error) {
	file, err := os.Open(fmt.Sprintf("%s/%s", e.RootDir, fname))
	if err != nil {
		return nil, fmt.Errorf("os: %s", err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)

	labels, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("labels: %s", err.Error())
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("records: %s", err.Error())
	}

	ret := make([]map[string]any, 0)

	for _, record := range records {

		retRecord := map[string]any{}

		for idx, val := range record {
			if idx >= len(labels) {
				return nil, fmt.Errorf("incorrect csv schema")
			}
			e.parseValue(retRecord, labels[idx], val)
		}

		ret = append(ret, retRecord)
	}

	return ret, nil
}

func (e *CSVEngine) List(path string) ([]string, error) {
	dir, err := os.Open(fmt.Sprintf("%s/%s", e.RootDir, path))
	if err != nil {
		return nil, fmt.Errorf("os: %s", err.Error())
	}
	defer dir.Close()

	files, err := dir.ReadDir(0)
	if err != nil {
		return nil, fmt.Errorf("os: %s", err.Error())
	}

	ret := make([]string, 0)
	for _, file := range files {
		ret = append(ret, fmt.Sprintf("%s/%s", path, file.Name()))
	}

	return ret, nil
}

func (e *CSVEngine) parseValue(record map[string]any, label string, value string) {

	valInt, err := strconv.Atoi(value)
	if err == nil {
		record[label] = valInt
		return
	}

	valFloat, err := strconv.ParseFloat(value, 64)
	if err == nil {
		record[label] = valFloat
		return
	}

	record[label] = value
}
