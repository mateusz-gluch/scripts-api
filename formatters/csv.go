package formatters

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/mitchellh/mapstructure"
)

type CSVFormatter[T any] struct{}

func (f *CSVFormatter[T]) Format(in []T) (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if len(in) == 0 {
		return "", nil
	}

	ret := make([]map[string]any, 0)
	for _, v := range in {
		retV := make(map[string]any)
		mapstructure.Decode(v, &retV)
		ret = append(ret, retV)
	}

	labels := make([]string, 0)
	for key := range ret[0] {
		labels = append(labels, key)
	}

	slices.Sort(labels)
	writer.Write(labels)

	for _, record := range ret {
		recordCsv := make([]string, len(labels))
		for k, v := range record {
			idx := slices.Index(labels, k)
			recordCsv[idx] = fmt.Sprintf("%v", v)
		}
		writer.Write(recordCsv)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("writer: %s", err.Error())
	}

	return buf.String(), nil
}
