package models

import (
	"encoding/json"
	"io"
)

// Преобразовывает тело запроса в виде JSON в соответствующую структуру
func JsonToStruct(structure any, req io.Reader) error {
	return json.NewDecoder(req).Decode(structure)
}

// Преобразовывает структуру в формат JSON
func StructToJson(structure any) ([]byte, error) {
	res, err := json.Marshal(structure)
	if err != nil {
		return nil, err
	}
	return res, nil
}
