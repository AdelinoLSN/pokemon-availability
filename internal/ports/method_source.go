package ports

import "github.com/AdelinoLSN/pokemon-availability/internal/domain"

type MethodSource interface {
	LoadMethods() ([]domain.Method, error)
}
