package ports

import "github.com/AdelinoLSN/pokemon-availability/internal/domain"

type GameSource interface {
	LoadGames() ([]domain.Game, error)
}
