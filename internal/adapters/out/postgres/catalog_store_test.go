package postgresadapter

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/importcatalog"
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestCatalogStoreUpsertsTheWholeCatalogInOneTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO games")).
		WithArgs("R", "Red", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO methods")).
		WithArgs("WILD", "Wild encounter").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO pokemon")).
		WithArgs(1, "Bulbasaur", "").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(25))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO pokemon_availability")).
		WithArgs(25, "R", "WILD", "").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta("REFRESH MATERIALIZED VIEW mv_pokemon_availability_details")).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = NewCatalogStore(db).Upsert(context.Background(), catalogFixture())

	if err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestCatalogStoreRollsBackWhenProjectionRefreshFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO games")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO methods")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO pokemon")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(25))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO pokemon_availability")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	refreshErr := errors.New("refresh failed")
	mock.ExpectExec(regexp.QuoteMeta("REFRESH MATERIALIZED VIEW mv_pokemon_availability_details")).
		WillReturnError(refreshErr)
	mock.ExpectRollback()

	err = NewCatalogStore(db).Upsert(context.Background(), catalogFixture())

	if !errors.Is(err, refreshErr) {
		t.Fatalf("Upsert() error = %v, want %v", err, refreshErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func catalogFixture() importcatalog.Catalog {
	return importcatalog.Catalog{
		Games:   []domain.Game{{Abbreviation: "R", Name: "Red", Generation: 1}},
		Methods: []domain.Method{{Key: "WILD", Description: "Wild encounter"}},
		Pokemons: []importcatalog.Pokemon{{
			Number: 1,
			Name:   "Bulbasaur",
			Availabilities: []importcatalog.Availability{{
				GameAbbreviation: "R",
				MethodKey:        "WILD",
			}},
		}},
	}
}
