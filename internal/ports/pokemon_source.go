package ports

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/dto"
)

type PokemonSource interface {
	LoadPokemonsJson() ([]dto.PokemonJson, error)
}
