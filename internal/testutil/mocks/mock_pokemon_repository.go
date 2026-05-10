package mocks

import (
	"errors"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

var _ ports.PokemonRepository = (*MockPokemonRepository)(nil)

type SaveResponse struct {
	Id  int
	Err error
}

type MockPokemonRepository struct {
	SaveResponses []SaveResponse
}

func (m *MockPokemonRepository) Save(_ domain.Pokemon) (int, error) {
	if len(m.SaveResponses) == 0 {
		return 0, errors.New("no mocked response configured")
	}

	response := m.SaveResponses[0]

	m.SaveResponses = m.SaveResponses[1:]

	return response.Id, response.Err
}
