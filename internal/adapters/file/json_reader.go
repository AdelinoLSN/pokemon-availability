package file

import (
	"encoding/json"
	"os"
)

type JSONReader struct{}

func NewJSONReader() *JSONReader {
	return &JSONReader{}
}

func (r *JSONReader) Read(path string, v any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(v)
}
