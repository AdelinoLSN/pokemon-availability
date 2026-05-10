package usecases

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

type LoadMethods struct {
	source ports.MethodSource
}

func NewLoadMethods(source ports.MethodSource) *LoadMethods {
	return &LoadMethods{
		source: source,
	}
}

func (u *LoadMethods) Execute() ([]domain.Method, error) {
	return u.source.LoadMethods()
}
