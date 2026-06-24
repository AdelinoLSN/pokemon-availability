package main

import (
	"errors"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/testutil"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRunExporter(t *testing.T) {
	t.Run("Success full flow", func(t *testing.T) {
		mockDB, mockSQL, err := sqlmock.New()
		assert.NoError(t, err)
		defer mockDB.Close()

		gameContent := `[{"abbreviation": "R", "name": "Red", "generation": 1}]`
		gamesFile := testutil.CreateTempFile(t, "", "games_*.json", gameContent)
		defer gamesFile.Close()

		columns := []string{"number", "name", "form", "game_abbreviation", "method_key", "note", "method_description", "game", "id"}
		mockSQL.ExpectQuery(`SELECT .*`).
			WillReturnRows(sqlmock.NewRows(columns).
				AddRow(
					25,           // number
					"Pikachu",    // name
					"",           // form
					"R",          // game_abbreviation
					"STARTER",    // method_key
					"Only one",   // note
					"Starter DB", // method_description
					"Red",        // game
					99,           // id
				))

		err = runExporter(mockDB, gamesFile.Name())

		assert.NoError(t, err)
		err = mockSQL.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Fails when loadGames fails", func(t *testing.T) {
		mockDB, _, _ := sqlmock.New()
		defer mockDB.Close()

		err := runExporter(mockDB, "/path/invalid/games.json")

		assert.Error(t, err)
	})

	t.Run("Fails when exportPokemonAvailabilityDetails fails", func(t *testing.T) {
		mockDB, mockSQL, err := sqlmock.New()
		assert.NoError(t, err)
		defer mockDB.Close()

		gameContent := `[{"abbreviation": "R", "name": "Red", "generation": 1}]`
		gamesFile := testutil.CreateTempFile(t, "", "games_*.json", gameContent)
		defer gamesFile.Close()

		mockSQL.ExpectQuery(`SELECT .*`).
			WillReturnError(errors.New("simulated db error"))

		err = runExporter(mockDB, gamesFile.Name())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "simulated db error")
	})
}

func TestLoadGames(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		gameContent := `[{"abbreviation": "R", "name": "Red", "generation": 1}]`
		gamesFile := testutil.CreateTempFile(t, "", "games_*.json", gameContent)
		defer gamesFile.Close()

		games, err := loadGames(gamesFile.Name())

		assert.NoError(t, err)
		assert.Len(t, games, 1)
		assert.Equal(t, "R", games[0].Abbreviation)
	})

	t.Run("Fails on invalid source", func(t *testing.T) {
		games, err := loadGames("/caminho/inexistente.json")

		assert.Error(t, err)
		assert.Nil(t, games)
	})
}

func TestExportPokemonAvailabilityDetails(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB, mockSQL, err := sqlmock.New()
		assert.NoError(t, err)
		defer mockDB.Close()

		games := []domain.Game{
			{Abbreviation: "R", Name: "Red", Generation: 1},
		}

		columns := []string{"number", "name", "form", "game_abbreviation", "method_key", "note", "method_description", "game", "id"}
		mockSQL.ExpectQuery(`SELECT .*`).
			WillReturnRows(sqlmock.NewRows(columns).
				AddRow(
					25,           // number
					"Pikachu",    // name
					"",           // form
					"R",          // game_abbreviation
					"STARTER",    // method_key
					"Only one",   // note
					"Starter DB", // method_description
					"Red",        // game
					99,           // id
				))

		err = exportPokemonAvailabilityDetails(mockDB, games)

		assert.NoError(t, err)
		err = mockSQL.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Fails on database error", func(t *testing.T) {
		mockDB, mockSQL, err := sqlmock.New()
		assert.NoError(t, err)
		defer mockDB.Close()

		games := []domain.Game{
			{Abbreviation: "R", Name: "Red", Generation: 1},
		}

		// Simula uma falha na query do banco de dados para forçar o retorno do erro
		mockSQL.ExpectQuery(`SELECT .*`).
			WillReturnError(errors.New("simulated database failure"))

		err = exportPokemonAvailabilityDetails(mockDB, games)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "simulated database failure")
	})
}
