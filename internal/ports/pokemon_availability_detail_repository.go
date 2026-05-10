package ports

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
)

type PokemonAvailabilityDetailRepository interface {
	RefreshMaterializedView() error
	LoadByGameAbbreviation(string) ([]models.PokemonAvailabilityDetail, error)
}
