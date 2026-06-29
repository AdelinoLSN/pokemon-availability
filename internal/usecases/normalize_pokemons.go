package usecases

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
)

type NormalizePokemon struct {
	source ports.PokemonSource
}

func NewNormalizePokemon(source ports.PokemonSource) *NormalizePokemon {
	return &NormalizePokemon{
		source: source,
	}
}

func (u *NormalizePokemon) Execute() ([]models.NormalizedPokemon, error) {
	return u.source.LoadNormalizedPokemons()
}
