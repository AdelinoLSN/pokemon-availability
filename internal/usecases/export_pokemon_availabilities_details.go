package usecases

import (
	"fmt"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
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
	path := u.buildExportPath(iterator, game)

	return u.exporter.ExportPokemonAvailabilityDetails(path, pokemonAvailabilityDetails)
}

func (u *ExportPokemonAvailabilityDetails) buildExportPath(iterator int, game domain.Game) string {
	filename := fmt.Sprintf("%03d_%s.csv", iterator, game.Abbreviation)

	return fmt.Sprintf(".outputs/%s", filename)
}
