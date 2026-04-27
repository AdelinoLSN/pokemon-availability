package ports

type JSONFileReader interface {
	Read(path string, v any) error
}

type CSVFileReader interface {
	ReadAll(path string) ([][]string, error)
}
