package file

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

type CSVHandler struct{}

func NewCSVHandler() *CSVHandler {
	return &CSVHandler{}
}

func (h *CSVHandler) ReadAll(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return csv.NewReader(file).ReadAll()
}

func (h *CSVHandler) WriteAll(path string, data [][]string) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.WriteAll(data)
}
