package ports

import "github.com/AdelinoLSN/pokemon-availability/internal/domain"

type MethodRepository interface {
	Save(domain.Method) error
	SaveAll([]domain.Method) error
}
