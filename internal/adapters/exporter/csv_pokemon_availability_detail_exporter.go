package exporter

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/infra/filesystem"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

var _ ports.PokemonAvailabilityDetailExporter = (*CsvPokemonAvailabilityDetailExporter)(nil)

type CsvPokemonAvailabilityDetailExporter struct{}

func NewCsvPokemonAvailabilityDetailExporter() *CsvPokemonAvailabilityDetailExporter {
	return &CsvPokemonAvailabilityDetailExporter{}
}

func (e *CsvPokemonAvailabilityDetailExporter) ExportPokemonAvailabilityDetails(path string, data [][]string) error {
	err := filesystem.WriteCSV(path, data)

	return err
}
