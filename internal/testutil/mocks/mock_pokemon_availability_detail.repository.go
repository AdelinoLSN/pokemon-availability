package mocks

import (
	"errors"

	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
)

var _ ports.PokemonAvailabilityDetailRepository = (*MockPokemonAvailabilityDetailRepository)(nil)

type RefreshMaterializedViewResponse struct {
	Err error
}

type LoadByGameAbbreviationResponse struct {
	Details []models.PokemonAvailabilityDetail
	Err     error
}

type MockPokemonAvailabilityDetailRepository struct {
	RefreshMaterializedViewResponses []RefreshMaterializedViewResponse
	LoadByGameAbbreviationResponses  []LoadByGameAbbreviationResponse

	RefreshMaterializedViewCalled bool
	LoadByGameAbbreviationCalls   []string
}

func (m *MockPokemonAvailabilityDetailRepository) RefreshMaterializedView() error {
	m.RefreshMaterializedViewCalled = true

	if len(m.RefreshMaterializedViewResponses) == 0 {
		return errors.New("no mocked response configured")
	}

	response := m.RefreshMaterializedViewResponses[0]

	m.RefreshMaterializedViewResponses =
		m.RefreshMaterializedViewResponses[1:]

	return response.Err
}

func (m *MockPokemonAvailabilityDetailRepository) LoadByGameAbbreviation(
	gameAbbreviation string,
) ([]models.PokemonAvailabilityDetail, error) {
	m.LoadByGameAbbreviationCalls = append(
		m.LoadByGameAbbreviationCalls,
		gameAbbreviation,
	)

	if len(m.LoadByGameAbbreviationResponses) == 0 {
		return nil, errors.New("no mocked response configured")
	}

	response := m.LoadByGameAbbreviationResponses[0]

	m.LoadByGameAbbreviationResponses =
		m.LoadByGameAbbreviationResponses[1:]

	return response.Details, response.Err
}
