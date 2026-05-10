package usecases

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/dto"
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
	pokemonsJson, err := u.source.LoadPokemonsJson()
	if err != nil {
		return nil, err
	}

	var normalizedPokemons []models.NormalizedPokemon

	for _, pokemonJson := range pokemonsJson {
		pokemonEntity := buildPokemonEntity(pokemonJson)

		var pokemonAvailabilities []domain.PokemonAvailability
		for _, pokemonAvailabilityJson := range pokemonJson.Availability {
			pokemonAvailabilityEntity := buildPokemonAvailabilityEntity(pokemonAvailabilityJson)
			pokemonAvailabilities = append(pokemonAvailabilities, pokemonAvailabilityEntity)
		}

		normalizedPokemons = append(
			normalizedPokemons,
			buildNormalizedPokemon(pokemonEntity, pokemonAvailabilities),
		)
	}

	return normalizedPokemons, nil
}

func buildPokemonEntity(pokemonJson dto.PokemonJson) domain.Pokemon {
	return domain.Pokemon{
		Number: pokemonJson.Number,
		Name:   pokemonJson.Name,
		Form:   pokemonJson.Form,
	}
}

func buildPokemonAvailabilityEntity(pokemonAvailabilityJson dto.PokemonAvailabilityJson) domain.PokemonAvailability {
	return domain.PokemonAvailability{
		GameAbbreviation: pokemonAvailabilityJson.Game,
		MethodKey:        pokemonAvailabilityJson.Method,
		Note:             pokemonAvailabilityJson.Notes,
	}
}

func buildNormalizedPokemon(
	pokemon domain.Pokemon,
	pokemonAvailabilities []domain.PokemonAvailability,
) models.NormalizedPokemon {
	return models.NormalizedPokemon{
		Pokemon:        pokemon,
		Availabilities: pokemonAvailabilities,
	}
}
