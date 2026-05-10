package usecases

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

type SaveGames struct {
	repository ports.GameRepository
}

func NewSaveGames(repository ports.GameRepository) *SaveGames {
	return &SaveGames{
		repository: repository,
	}
}

func (u *SaveGames) Execute(games []domain.Game) error {
	return u.repository.SaveAll(games)
}
