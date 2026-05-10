package ports

import "github.com/AdelinoLSN/pokemon-availability/internal/domain"

type GameRepository interface {
	GetAll() ([]domain.Game, error)
	Save(domain.Game) error
	SaveAll([]domain.Game) error
}
