package source

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/infra/filesystem"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

var _ ports.MethodSource = (*JsonMethodSource)(nil)

type JsonMethodSource struct {
	filepath string
}

func NewJsonMethodSource(filepath string) *JsonMethodSource {
	return &JsonMethodSource{
		filepath: filepath,
	}
}

func (s *JsonMethodSource) LoadMethods() ([]domain.Method, error) {
	var methodMap domain.MethodMap

	if err := filesystem.ReadJson(s.filepath, &methodMap); err != nil {
		return nil, err
	}

	var methods []domain.Method

	for k, d := range methodMap {
		methods = append(methods, domain.Method{
			Key:         k,
			Description: d,
		})
	}

	return methods, nil
}
