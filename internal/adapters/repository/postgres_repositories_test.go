package repository

import (
	"errors"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresGameRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT abbreviation, name, generation`).
		WillReturnRows(sqlmock.NewRows([]string{"abbreviation", "name", "generation"}).
			AddRow("R", "Red", 1).
			AddRow("B", "Blue", 1))

	games, err := NewPostgresGameRepository(db).GetAll()

	require.NoError(t, err)
	assert.Equal(t, []domain.Game{
		{Abbreviation: "R", Name: "Red", Generation: 1},
		{Abbreviation: "B", Name: "Blue", Generation: 1},
	}, games)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresGameRepository_GetAll_ReturnsQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	expectedErr := errors.New("query failed")
	mock.ExpectQuery(`SELECT abbreviation, name, generation`).WillReturnError(expectedErr)

	games, err := NewPostgresGameRepository(db).GetAll()

	assert.Nil(t, games)
	assert.ErrorIs(t, err, expectedErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresGameRepository_GetAll_ReturnsScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT abbreviation, name, generation`).
		WillReturnRows(sqlmock.NewRows([]string{"abbreviation", "name", "generation"}).
			AddRow("R", "Red", "invalid"))

	games, err := NewPostgresGameRepository(db).GetAll()

	assert.Nil(t, games)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresGameRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`INSERT INTO .* \(abbreviation, name, generation\)`).
		WithArgs("R", "Red", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = NewPostgresGameRepository(db).Save(domain.Game{Abbreviation: "R", Name: "Red", Generation: 1})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresGameRepository_SaveAll_ReturnsSaveError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	expectedErr := errors.New("save failed")
	mock.ExpectExec(`INSERT INTO .* \(abbreviation, name, generation\)`).
		WithArgs("R", "Red", 1).
		WillReturnError(expectedErr)

	err = NewPostgresGameRepository(db).SaveAll([]domain.Game{{Abbreviation: "R", Name: "Red", Generation: 1}})

	assert.ErrorIs(t, err, expectedErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresMethodRepository_SaveAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`INSERT INTO .* \(key, description\)`).
		WithArgs("WILD", "Wild encounter").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = NewPostgresMethodRepository(db).SaveAll([]domain.Method{{Key: "WILD", Description: "Wild encounter"}})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresMethodRepository_SaveAll_ReturnsSaveError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	expectedErr := errors.New("save failed")
	mock.ExpectExec(`INSERT INTO .* \(key, description\)`).
		WithArgs("WILD", "Wild encounter").
		WillReturnError(expectedErr)

	err = NewPostgresMethodRepository(db).SaveAll([]domain.Method{{Key: "WILD", Description: "Wild encounter"}})

	assert.ErrorIs(t, err, expectedErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`INSERT INTO .* \(number, name, form\)`).
		WithArgs(25, "Pikachu", "").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(99))

	id, err := NewPostgresPokemonRepository(db).Save(domain.Pokemon{Number: 25, Name: "Pikachu", Form: ""})

	require.NoError(t, err)
	assert.Equal(t, 99, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonRepository_Save_ReturnsQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	expectedErr := errors.New("save failed")
	mock.ExpectQuery(`INSERT INTO .* \(number, name, form\)`).
		WithArgs(25, "Pikachu", "").
		WillReturnError(expectedErr)

	id, err := NewPostgresPokemonRepository(db).Save(domain.Pokemon{Number: 25, Name: "Pikachu", Form: ""})

	assert.Zero(t, id)
	assert.ErrorIs(t, err, expectedErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonAvailabilityRepository_SaveAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`INSERT INTO .* \(pokemon_id, game_abbreviation, method_key, note\) VALUES \(\$1, \$2, \$3, \$4\),\(\$5, \$6, \$7, \$8\)`).
		WithArgs(99, "R", "WILD", "", 99, "B", "GIFT", "NPC").
		WillReturnResult(sqlmock.NewResult(0, 2))

	err = NewPostgresPokemonAvailabilityRepository(db).SaveAll([]domain.PokemonAvailability{
		{PokemonId: 99, GameAbbreviation: "R", MethodKey: "WILD", Note: ""},
		{PokemonId: 99, GameAbbreviation: "B", MethodKey: "GIFT", Note: "NPC"},
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonAvailabilityRepository_SaveAll_EmptySliceDoesNothing(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	err = NewPostgresPokemonAvailabilityRepository(db).SaveAll(nil)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonAvailabilityRepository_SaveAll_ReturnsExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	expectedErr := errors.New("save failed")
	mock.ExpectExec(`INSERT INTO .* \(pokemon_id, game_abbreviation, method_key, note\)`).
		WithArgs(99, "R", "WILD", "").
		WillReturnError(expectedErr)

	err = NewPostgresPokemonAvailabilityRepository(db).SaveAll([]domain.PokemonAvailability{
		{PokemonId: 99, GameAbbreviation: "R", MethodKey: "WILD", Note: ""},
	})

	assert.ErrorIs(t, err, expectedErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonAvailabilityDetailRepository_RefreshMaterializedView(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`REFRESH MATERIALIZED VIEW .*`).WillReturnResult(sqlmock.NewResult(0, 0))

	err = NewPostgresPokemonAvailabilityDetailRepository(db).RefreshMaterializedView()

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonAvailabilityDetailRepository_RefreshMaterializedView_ReturnsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	expectedErr := errors.New("refresh failed")
	mock.ExpectExec(`REFRESH MATERIALIZED VIEW .*`).WillReturnError(expectedErr)

	err = NewPostgresPokemonAvailabilityDetailRepository(db).RefreshMaterializedView()

	assert.ErrorIs(t, err, expectedErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonAvailabilityDetailRepository_LoadByGameAbbreviation(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	columns := []string{"number", "name", "form", "game_abbreviation", "game", "method_key", "method_description", "note", "id"}
	mock.ExpectQuery(`SELECT .* FROM .* WHERE game_abbreviation = \$1`).
		WithArgs("R").
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(25, "Pikachu", "", "R", "Red", "WILD", "Wild encounter", "", 99))

	details, err := NewPostgresPokemonAvailabilityDetailRepository(db).LoadByGameAbbreviation("R")

	require.NoError(t, err)
	require.Len(t, details, 1)
	assert.Equal(t, 25, details[0].Number)
	assert.Equal(t, "Pikachu", details[0].Name)
	assert.Equal(t, "R", details[0].GameAbbreviation)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonAvailabilityDetailRepository_LoadByGameAbbreviation_ReturnsQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	expectedErr := errors.New("query failed")
	mock.ExpectQuery(`SELECT .* FROM .* WHERE game_abbreviation = \$1`).
		WithArgs("R").
		WillReturnError(expectedErr)

	details, err := NewPostgresPokemonAvailabilityDetailRepository(db).LoadByGameAbbreviation("R")

	assert.Nil(t, details)
	assert.ErrorIs(t, err, expectedErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresPokemonAvailabilityDetailRepository_LoadByGameAbbreviation_ReturnsScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	columns := []string{"number", "name", "form", "game_abbreviation", "game", "method_key", "method_description", "note", "id"}
	mock.ExpectQuery(`SELECT .* FROM .* WHERE game_abbreviation = \$1`).
		WithArgs("R").
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow("invalid", "Pikachu", "", "R", "Red", "WILD", "Wild encounter", "", 99))

	details, err := NewPostgresPokemonAvailabilityDetailRepository(db).LoadByGameAbbreviation("R")

	assert.Nil(t, details)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
