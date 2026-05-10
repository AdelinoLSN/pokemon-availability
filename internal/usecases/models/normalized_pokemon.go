package models

import "github.com/AdelinoLSN/pokemon-availability/internal/domain"

type NormalizedPokemon struct {
	Pokemon        domain.Pokemon
	Availabilities []domain.PokemonAvailability
}
