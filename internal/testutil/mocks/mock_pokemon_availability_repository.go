package mocks

import (
	"errors"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

var _ ports.PokemonAvailabilityRepository = (*MockPokemonAvailabilityRepository)(nil)

type SaveAllResponse struct {
	Err error
}

type MockPokemonAvailabilityRepository struct {
	SaveAllResponses    []SaveAllResponse
	SavedAvailabilities [][]domain.PokemonAvailability
}

func (m *MockPokemonAvailabilityRepository) SaveAll(availabilities []domain.PokemonAvailability) error {
	m.SavedAvailabilities = append(
		m.SavedAvailabilities,
		availabilities,
	)

	if len(m.SaveAllResponses) == 0 {
		return errors.New("no mocked response configured")
	}

	response := m.SaveAllResponses[0]

	m.SaveAllResponses = m.SaveAllResponses[1:]

	return response.Err
}
