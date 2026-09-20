package importcatalog

import (
	"context"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
)

type Command struct{}

type InputPort interface {
	Execute(context.Context, Command) error
}

type Source interface {
	Load(context.Context) (Catalog, error)
}

type Store interface {
	Upsert(context.Context, Catalog) error
}

type Catalog struct {
	Games    []domain.Game
	Methods  []domain.Method
	Pokemons []Pokemon
}

type Pokemon struct {
	Number         int
	Name           string
	Form           string
	Availabilities []Availability
}

type Availability struct {
	GameAbbreviation string
	MethodKey        string
	Note             string
}

type UseCase struct {
	source Source
	store  Store
}

func New(source Source, store Store) *UseCase {
	return &UseCase{source: source, store: store}
}

func (u *UseCase) Execute(ctx context.Context, _ Command) error {
	catalog, err := u.source.Load(ctx)
	if err != nil {
		return err
	}

	return u.store.Upsert(ctx, catalog)
}
