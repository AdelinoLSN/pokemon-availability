package importcatalog_test

import (
	"context"
	"errors"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/importcatalog"
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
)

func TestUseCaseImportsTheCatalogLoadedByTheSource(t *testing.T) {
	catalog := importcatalog.Catalog{
		Games:   []domain.Game{{Abbreviation: "R", Name: "Red", Generation: 1}},
		Methods: []domain.Method{{Key: "WILD", Description: "Wild encounter"}},
	}
	source := &sourceStub{catalog: catalog}
	store := &storeSpy{}

	err := importcatalog.New(source, store).Execute(context.Background(), importcatalog.Command{})

	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !store.called {
		t.Fatal("store was not called")
	}
	if store.catalog.Games[0].Abbreviation != "R" {
		t.Fatalf("stored game = %q, want R", store.catalog.Games[0].Abbreviation)
	}
}

func TestUseCaseDoesNotStoreWhenSourceFails(t *testing.T) {
	source := &sourceStub{err: errors.New("source unavailable")}
	store := &storeSpy{}

	err := importcatalog.New(source, store).Execute(context.Background(), importcatalog.Command{})

	if !errors.Is(err, source.err) {
		t.Fatalf("Execute() error = %v, want %v", err, source.err)
	}
	if store.called {
		t.Fatal("store must not be called after a source failure")
	}
}

func TestUseCaseReturnsStoreFailure(t *testing.T) {
	store := &storeSpy{err: errors.New("database unavailable")}
	useCase := importcatalog.New(&sourceStub{}, store)

	err := useCase.Execute(context.Background(), importcatalog.Command{})

	if !errors.Is(err, store.err) {
		t.Fatalf("Execute() error = %v, want %v", err, store.err)
	}
}

type sourceStub struct {
	catalog importcatalog.Catalog
	err     error
}

func (s *sourceStub) Load(context.Context) (importcatalog.Catalog, error) {
	return s.catalog, s.err
}

type storeSpy struct {
	called  bool
	catalog importcatalog.Catalog
	err     error
}

func (s *storeSpy) Upsert(_ context.Context, catalog importcatalog.Catalog) error {
	s.called = true
	s.catalog = catalog
	return s.err
}
