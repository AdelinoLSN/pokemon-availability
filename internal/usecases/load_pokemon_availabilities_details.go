package usecases

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
)

type LoadPokemonAvailabilityDetails struct {
	repository ports.PokemonAvailabilityDetailRepository
}

func NewLoadPokemonAvailabilityDetails(
	repository ports.PokemonAvailabilityDetailRepository,
) *LoadPokemonAvailabilityDetails {
	return &LoadPokemonAvailabilityDetails{
		repository: repository,
	}
}

func (u *LoadPokemonAvailabilityDetails) Execute(gameAbbreviation string) ([]models.PokemonAvailabilityDetail, error) {
	pokemonAvailabilityDetails, err := u.repository.LoadByGameAbbreviation(gameAbbreviation)

	if err != nil {
		return nil, err
	}

	return pokemonAvailabilityDetails, nil
}
