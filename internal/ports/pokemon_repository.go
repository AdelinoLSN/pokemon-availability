package ports

import "github.com/AdelinoLSN/pokemon-availability/internal/domain"

type PokemonRepository interface {
	Save(domain.Pokemon) (int, error)
}
