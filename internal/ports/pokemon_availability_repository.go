package ports

import "github.com/AdelinoLSN/pokemon-availability/internal/domain"

type PokemonAvailabilityRepository interface {
	SaveAll([]domain.PokemonAvailability) error
}
