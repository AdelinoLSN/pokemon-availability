package mocks

import (
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockPokemonRepository_Save(t *testing.T) {
	repository := &MockPokemonRepository{
		SaveResponses: []SaveResponse{{Id: 25}},
	}

	id, err := repository.Save(domain.Pokemon{Number: 25})

	require.NoError(t, err)
	assert.Equal(t, 25, id)
}

func TestMockPokemonRepository_Save_ReturnsErrorWhenResponseIsMissing(t *testing.T) {
	id, err := (&MockPokemonRepository{}).Save(domain.Pokemon{})

	assert.Zero(t, id)
	assert.Error(t, err)
}

func TestMockPokemonAvailabilityRepository_SaveAll(t *testing.T) {
	availabilities := []domain.PokemonAvailability{{PokemonId: 25}}
	repository := &MockPokemonAvailabilityRepository{
		SaveAllResponses: []SaveAllResponse{{Err: nil}},
	}

	err := repository.SaveAll(availabilities)

	require.NoError(t, err)
	assert.Equal(t, [][]domain.PokemonAvailability{availabilities}, repository.SavedAvailabilities)
}

func TestMockPokemonAvailabilityRepository_SaveAll_ReturnsErrorWhenResponseIsMissing(t *testing.T) {
	err := (&MockPokemonAvailabilityRepository{}).SaveAll(nil)

	assert.Error(t, err)
}

func TestMockPokemonAvailabilityDetailRepository_RefreshMaterializedView(t *testing.T) {
	repository := &MockPokemonAvailabilityDetailRepository{
		RefreshMaterializedViewResponses: []RefreshMaterializedViewResponse{{Err: nil}},
	}

	err := repository.RefreshMaterializedView()

	require.NoError(t, err)
	assert.True(t, repository.RefreshMaterializedViewCalled)
}

func TestMockPokemonAvailabilityDetailRepository_RefreshMaterializedView_ReturnsErrorWhenResponseIsMissing(t *testing.T) {
	err := (&MockPokemonAvailabilityDetailRepository{}).RefreshMaterializedView()

	assert.Error(t, err)
}

func TestMockPokemonAvailabilityDetailRepository_LoadByGameAbbreviation(t *testing.T) {
	details := []models.PokemonAvailabilityDetail{{Number: 25}}
	repository := &MockPokemonAvailabilityDetailRepository{
		LoadByGameAbbreviationResponses: []LoadByGameAbbreviationResponse{{Details: details}},
	}

	got, err := repository.LoadByGameAbbreviation("R")

	require.NoError(t, err)
	assert.Equal(t, details, got)
	assert.Equal(t, []string{"R"}, repository.LoadByGameAbbreviationCalls)
}

func TestMockPokemonAvailabilityDetailRepository_LoadByGameAbbreviation_ReturnsErrorWhenResponseIsMissing(t *testing.T) {
	details, err := (&MockPokemonAvailabilityDetailRepository{}).LoadByGameAbbreviation("R")

	assert.Nil(t, details)
	assert.Error(t, err)
}
