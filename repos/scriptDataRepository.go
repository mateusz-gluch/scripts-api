package repos

import (
	"fmt"
	"scripts-api/dataengine"
	"scripts-api/model"
	"strings"
	"time"
)

type ScriptDataRepository struct {
	Engine dataengine.FileEngine
}

// READ ONLY
func (repo *ScriptDataRepository) SelectSpec(key string, spec *model.ScriptSpec) (*[]map[string]any, error) {

	dirpath := fmt.Sprintf("%s/organizationId=%d", key, spec.Organization)
	files, err := repo.Engine.List(dirpath)
	if err != nil {
		return nil, fmt.Errorf("path: %s", err.Error())
	}

	var filesToRead []string

	if len(spec.Assets) == 0 {
		filesToRead = files
	} else {
		filesToRead = repo.filterFiles(files, spec.Assets)
	}

	ret := make([]map[string]any, 0)

	for _, file := range filesToRead {
		contents, err := repo.Engine.Read(file)
		if err != nil {
			return nil, fmt.Errorf("engine: %s", err.Error())
		}
		ret = append(ret, contents...)
	}

	ret, err = repo.filterTimestamp(ret, spec)
	if err != nil {
		return nil, fmt.Errorf("timestamp: %s", err.Error())
	}

	ret, err = repo.filterOptions(ret, spec)
	if err != nil {
		return nil, fmt.Errorf("options: %s", err.Error())
	}

	// fmt.Println(ret)

	return &ret, nil
}

func (repo *ScriptDataRepository) filterFiles(files []string, assets []int) []string {
	ret := make([]string, 0)

	for _, file := range files {
		for _, assetId := range assets {
			assetIdentifier := fmt.Sprintf("assetId=%d", assetId)
			if strings.Contains(file, assetIdentifier) {
				ret = append(ret, file)
			}
		}
	}

	return ret
}

func (repo *ScriptDataRepository) filterOptions(
	records []map[string]any,
	spec *model.ScriptSpec,
) ([]map[string]any, error) {

	ret := make([]map[string]any, 0, len(records))
	ret = append(ret, records...)

	var err error

	for opt, val := range spec.Options {
		ret, err = repo.filterField(ret, opt, val)
		if err != nil {
			return nil, fmt.Errorf("field: %s", err.Error())
		}
	}

	return ret, nil
}

func (repo *ScriptDataRepository) filterField(
	records []map[string]any,
	field string,
	vals []string,
) ([]map[string]any, error) {
	ret := make([]map[string]any, 0)

	for _, record := range records {

		val, ok := record[field].(string)
		if !ok {
			return nil, fmt.Errorf("field not provided in schema - spec: %s, %+v", field, vals)
		}

		for _, v := range vals {
			if v == val {
				ret = append(ret, record)
			}
		}

	}

	return ret, nil

}

func (repo *ScriptDataRepository) filterTimestamp(
	records []map[string]any,
	spec *model.ScriptSpec,
) ([]map[string]any, error) {

	ret := make([]map[string]any, 0)

	for _, record := range records {
		ts, ok := record["workingDay"]
		if !ok {
			return nil, fmt.Errorf("timestamp not provided in schema - spec: %+v, record: %+v", spec, record)
		}

		tz, ok := record["timezone"]
		if !ok {
			return nil, fmt.Errorf("timezone not provided in schema - spec: %+v", spec)
		}

		loc, err := time.LoadLocation(tz.(string))
		if err != nil {
			return nil, fmt.Errorf("timezone: %s", err.Error())
		}
		recordStr := strings.Replace(ts.(string), "T", " ", 1)
		recordTsEnd, err := time.ParseInLocation("2006-01-02", recordStr, loc)
		if err != nil {
			return nil, fmt.Errorf("timestamp provided in incorrect format: %s spec: %+v", err.Error(), spec)
		}

		recordTsStart := recordTsEnd.Add(time.Duration(24 * time.Hour))

		if recordTsStart.After(spec.StartTs) && recordTsEnd.Before(spec.EndTs) {
			ret = append(ret, record)
		}
	}

	return ret, nil
}
