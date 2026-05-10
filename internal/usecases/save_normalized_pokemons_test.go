package usecases_test

import (
	"errors"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/testutil/mocks"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
)

func TestSaveNormalizedPokemons_ShouldReturnError_WhenSaveAllFails(t *testing.T) {
	expectedErr := errors.New("failed to save availabilities")

	pokemonRepository := &mocks.MockPokemonRepository{
		SaveResponses: []mocks.SaveResponse{
			{
				Id: 25,
			},
		},
	}

	pokemonAvailabilityRepository := &mocks.MockPokemonAvailabilityRepository{
		SaveAllResponses: []mocks.SaveAllResponse{
			{
				Err: expectedErr,
			},
		},
	}

	pokemonAvailabilityDetailRepository :=
		&mocks.MockPokemonAvailabilityDetailRepository{
			RefreshMaterializedViewResponses: []mocks.RefreshMaterializedViewResponse{
				{
					Err: nil,
				},
			},
		}

	usecase := usecases.NewSaveNormalizedPokemons(
		pokemonRepository,
		pokemonAvailabilityRepository,
		pokemonAvailabilityDetailRepository,
	)

	normalizedPokemons := []models.NormalizedPokemon{
		{
			Pokemon: domain.Pokemon{
				Number: 25,
				Name:   "Pikachu",
			},
			Availabilities: []domain.PokemonAvailability{
				{
					GameAbbreviation: "RED",
					MethodKey:        "STARTER",
					Note:             "Only one",
				},
			},
		},
	}

	err := usecase.Execute(normalizedPokemons)

	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v but got %v", expectedErr, err)
	}

	if pokemonAvailabilityDetailRepository.RefreshMaterializedViewCalled {
		t.Fatalf("expected RefreshMaterializedView to not be called")
	}
}
