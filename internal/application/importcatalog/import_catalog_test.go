package importcatalog_test

import (
	"context"
	"errors"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/importcatalog"
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
)

func TestUseCaseLoadsAndStoresCatalog(t *testing.T) {
	source := &sourceStub{catalog: importcatalog.Catalog{
		Games: []domain.Game{{Abbreviation: "R"}},
	}}
	store := &storeSpy{}

	err := importcatalog.New(source, store).Execute(context.Background(), importcatalog.Command{})

	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !store.called || store.catalog.Games[0].Abbreviation != "R" {
		t.Fatalf("store = %#v", store)
	}
}

func TestUseCaseDoesNotStoreWhenSourceFails(t *testing.T) {
	expected := errors.New("source unavailable")
	store := &storeSpy{}

	err := importcatalog.New(&sourceStub{err: expected}, store).Execute(context.Background(), importcatalog.Command{})

	if !errors.Is(err, expected) {
		t.Fatalf("Execute() error = %v, want %v", err, expected)
	}
	if store.called {
		t.Fatal("store must not be called")
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
}

func (s *storeSpy) Upsert(_ context.Context, catalog importcatalog.Catalog) error {
	s.called = true
	s.catalog = catalog
	return nil
}
