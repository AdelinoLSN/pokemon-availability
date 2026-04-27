package file

import (
	"encoding/csv"
	"os"
)

type CSVReader struct{}

func NewCSVReader() *CSVReader {
	return &CSVReader{}
}

func (r *CSVReader) ReadAll(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return csv.NewReader(file).ReadAll()
}
