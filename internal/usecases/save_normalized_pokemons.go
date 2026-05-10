package usecases

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
)

type SaveNormalizedPokemons struct {
	pokemonRepository                   ports.PokemonRepository
	pokemonAvailabilityRepository       ports.PokemonAvailabilityRepository
	pokemonAvailabilityDetailRepository ports.PokemonAvailabilityDetailRepository
}

func NewSaveNormalizedPokemons(
	pokemonRepository ports.PokemonRepository,
	pokemonAvailabilityRepository ports.PokemonAvailabilityRepository,
	pokemonAvailabilityDetailRepository ports.PokemonAvailabilityDetailRepository,
) *SaveNormalizedPokemons {
	return &SaveNormalizedPokemons{
		pokemonRepository:                   pokemonRepository,
		pokemonAvailabilityRepository:       pokemonAvailabilityRepository,
		pokemonAvailabilityDetailRepository: pokemonAvailabilityDetailRepository,
	}
}

func (u *SaveNormalizedPokemons) Execute(normalizedPokemons []models.NormalizedPokemon) error {
	for _, normalizedPokemon := range normalizedPokemons {
		pokemonId, err := u.pokemonRepository.Save(normalizedPokemon.Pokemon)
		if err != nil {
			return err
		}

		for i := range normalizedPokemon.Availabilities {
			normalizedPokemon.Availabilities[i].PokemonId = pokemonId
		}

		if err := u.pokemonAvailabilityRepository.SaveAll(normalizedPokemon.Availabilities); err != nil {
			return err
		}
	}

	u.pokemonAvailabilityDetailRepository.RefreshMaterializedView()

	return nil
}
