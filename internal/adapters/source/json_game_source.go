package source

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/infra/filesystem"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

var _ ports.GameSource = (*JsonGameSource)(nil)

type JsonGameSource struct {
	filepath string
}

func NewJsonGameSource(filepath string) *JsonGameSource {
	return &JsonGameSource{
		filepath: filepath,
	}
}

func (s *JsonGameSource) LoadGames() ([]domain.Game, error) {
	var games []domain.Game

	err := filesystem.ReadJson(s.filepath, &games)

	return games, err
}
