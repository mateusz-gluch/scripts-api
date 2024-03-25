package formatters

import "encoding/json"

type JSONFormatter[T any] struct {
}

func (f *JSONFormatter[T]) Format(in []T) (string, error) {
	ret, err := json.Marshal(in)
	return string(ret), err
}
