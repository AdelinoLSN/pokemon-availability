package usecases

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

type LoadGames struct {
	source ports.GameSource
}

func NewLoadGames(source ports.GameSource) *LoadGames {
	return &LoadGames{
		source: source,
	}
}

func (u *LoadGames) Execute() ([]domain.Game, error) {
	return u.source.LoadGames()
}
