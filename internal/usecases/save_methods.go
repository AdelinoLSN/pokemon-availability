package usecases

import (
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

type SaveMethods struct {
	repository ports.MethodRepository
}

func NewSaveMethods(repository ports.MethodRepository) *SaveMethods {
	return &SaveMethods{
		repository: repository,
	}
}

func (u *SaveMethods) Execute(methods []domain.Method) error {
	return u.repository.SaveAll(methods)
}
