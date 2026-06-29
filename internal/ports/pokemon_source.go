package ports

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
)

type PokemonSource interface {
	LoadNormalizedPokemons() ([]models.NormalizedPokemon, error)
}
