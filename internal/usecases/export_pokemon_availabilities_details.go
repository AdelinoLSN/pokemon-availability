package usecases

import (
	"fmt"
	"strings"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/dto"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
)

type ExportPokemonAvailabilityDetails struct {
	exporter ports.PokemonAvailabilityDetailExporter
}

func NewExportPokemonAvailabilityDetails(
	exporter ports.PokemonAvailabilityDetailExporter,
) *ExportPokemonAvailabilityDetails {
	return &ExportPokemonAvailabilityDetails{
		exporter: exporter,
	}
}

func (u *ExportPokemonAvailabilityDetails) Execute(
	iterator int,
	game domain.Game,
	pokemonAvailabilityDetails []models.PokemonAvailabilityDetail,
) error {
	csvRows := u.buildCSVRows(pokemonAvailabilityDetails)

	csvData := u.toCSVData(csvRows)

	path := u.buildExportPath(iterator, game)

	return u.exporter.ExportPokemonAvailabilityDetails(path, csvData)
}

func (u *ExportPokemonAvailabilityDetails) buildExportPath(iterator int, game domain.Game) string {
	filename := fmt.Sprintf("%03d_%s.csv", iterator, game.Abbreviation)

	return fmt.Sprintf(".outputs/%s", filename)
}

func (u *ExportPokemonAvailabilityDetails) buildCSVRows(rows []models.PokemonAvailabilityDetail) []dto.PokemonCSVRow {
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

		currentNotes = append(currentNotes, row.Note)
	}

	appendCurrentPokemon()

	return csvRows
}

func (u *ExportPokemonAvailabilityDetails) toCSVData(rows []dto.PokemonCSVRow) [][]string {
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
