package exporter

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCsvPokemonAvailabilityDetailExporter_ExportPokemonAvailabilityDetails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "exports", "pokemon.csv")
	details := []models.PokemonAvailabilityDetail{
		{Number: 25, Name: "Pikachu", Form: "", MethodKey: "WILD", Note: ""},
		{Number: 25, Name: "Pikachu", Form: "", MethodKey: "GIFT", Note: "NPC, city"},
		{Number: 26, Name: "Raichu", Form: "Alola", MethodKey: "TRADE", Note: "Island"},
	}

	err := NewCsvPokemonAvailabilityDetailExporter().ExportPokemonAvailabilityDetails(path, details)

	require.NoError(t, err)
	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	require.NoError(t, err)
	assert.Equal(t, [][]string{
		{"Number", "Name", "Form", "Methods", "Notes"},
		{"0025", "Pikachu", "", "WILD|GIFT", "|[GIFT]: NPC; city"},
		{"0026", "Raichu", "Alola", "TRADE", "[TRADE]: Island"},
	}, records)
}

func TestCsvPokemonAvailabilityDetailExporter_ExportPokemonAvailabilityDetails_EmptyDetails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pokemon.csv")

	err := NewCsvPokemonAvailabilityDetailExporter().ExportPokemonAvailabilityDetails(path, nil)

	require.NoError(t, err)
	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"Number", "Name", "Form", "Methods", "Notes"}}, records)
}
