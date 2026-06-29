package exporter

import (
	"fmt"
	"strings"

	"github.com/AdelinoLSN/pokemon-availability/internal/dto"
	"github.com/AdelinoLSN/pokemon-availability/internal/infra/filesystem"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
)

var _ ports.PokemonAvailabilityDetailExporter = (*CsvPokemonAvailabilityDetailExporter)(nil)

type CsvPokemonAvailabilityDetailExporter struct{}

func NewCsvPokemonAvailabilityDetailExporter() *CsvPokemonAvailabilityDetailExporter {
	return &CsvPokemonAvailabilityDetailExporter{}
}

func (e *CsvPokemonAvailabilityDetailExporter) ExportPokemonAvailabilityDetails(
	path string,
	details []models.PokemonAvailabilityDetail,
) error {
	data := e.toCSVData(e.buildCSVRows(details))

	err := filesystem.WriteCSV(path, data)

	return err
}

func (e *CsvPokemonAvailabilityDetailExporter) buildCSVRows(
	rows []models.PokemonAvailabilityDetail,
) []dto.PokemonCSVRow {
	if len(rows) == 0 {
		return []dto.PokemonCSVRow{}
	}

	var csvRows []dto.PokemonCSVRow

	var (
		currentNumber  int
		currentName    string
		currentForm    string
		currentMethods []string
		currentNotes   []string
	)

	appendCurrentPokemon := func() {
		methods := strings.Join(currentMethods, "|")

		notes := strings.Join(currentNotes, "|")
		notes = strings.ReplaceAll(notes, ",", ";")

		csvRows = append(csvRows, dto.PokemonCSVRow{
			Number:  fmt.Sprintf("%04d", currentNumber),
			Name:    currentName,
			Form:    currentForm,
			Methods: methods,
			Notes:   notes,
		})
	}

	for i, row := range rows {
		isNewPokemon := i == 0 || row.Number != currentNumber || row.Form != currentForm

		if isNewPokemon {
			if i > 0 {
				appendCurrentPokemon()
			}

			currentNumber = row.Number
			currentName = row.Name
			currentForm = row.Form

			currentMethods = []string{}
			currentNotes = []string{}
		}

		currentMethods = append(currentMethods, row.MethodKey)
		rowNote := row.Note
		if rowNote != "" {
			rowNote = fmt.Sprintf("[%s]: %s", row.MethodKey, row.Note)
		}

		currentNotes = append(currentNotes, rowNote)
	}

	appendCurrentPokemon()

	return csvRows
}

func (e *CsvPokemonAvailabilityDetailExporter) toCSVData(rows []dto.PokemonCSVRow) [][]string {
	data := [][]string{
		{
			"Number",
			"Name",
			"Form",
			"Methods",
			"Notes",
		},
	}

	for _, row := range rows {
		data = append(data, []string{
			row.Number,
			row.Name,
			row.Form,
			row.Methods,
			row.Notes,
		})
	}

	return data
}
