package ports

import "github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"

type PokemonAvailabilityDetailExporter interface {
	ExportPokemonAvailabilityDetails(path string, details []models.PokemonAvailabilityDetail) error
}
