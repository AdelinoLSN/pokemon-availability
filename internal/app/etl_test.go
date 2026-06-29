package app

import (
	"errors"
	"os"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/testutil"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRunETL(t *testing.T) {
	t.Run("Success full flow", func(t *testing.T) {
		gameContent := `[{"abbreviation": "R", "name": "Red", "generation": 1}]`
		gamesFile := testutil.CreateTempFile(t, "", "games_*.json", gameContent)
		defer gamesFile.Close()

		methodContent := `{"WILD": "Wild encounter"}`
		methodsFile := testutil.CreateTempFile(t, "", "methods_*.json", methodContent)
		defer methodsFile.Close()

		pokemonsDir := testutil.CreateTempDir(t, "", "pokemons_dir_*")
		defer os.RemoveAll(pokemonsDir)

		pokemonContent := `[{"number": 1, "name": "Bulbasaur", "form": "", "availability": [{"game": "R", "method": "WILD", "notes": ""}]}]`
		testutil.CreateFile(t, pokemonsDir, "1.json", pokemonContent)

		mockDB, mockSQL, err := sqlmock.New()
		assert.NoError(t, err)
		defer mockDB.Close()

		mockSQL.ExpectExec(`INSERT INTO .* \(abbreviation, name, generation\)`).
			WithArgs("R", "Red", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockSQL.ExpectExec(`INSERT INTO .* \(key, description\)`).
			WithArgs("WILD", "Wild encounter").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockSQL.ExpectQuery(`INSERT INTO .* \(number, name, form\)`).
			WithArgs(1, "Bulbasaur", "").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(99))

		mockSQL.ExpectExec(`INSERT INTO .* \(pokemon_id, game_abbreviation, method_key, note\) VALUES \(\$1, \$2, \$3, \$4\)`).
			WithArgs(99, "R", "WILD", "").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockSQL.ExpectExec(`REFRESH MATERIALIZED VIEW .*`).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err = RunETL(mockDB, gamesFile.Name(), methodsFile.Name(), pokemonsDir)

		assert.NoError(t, err)
		err = mockSQL.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Fails on step 1 persistGames", func(t *testing.T) {
		mockDB, _, _ := sqlmock.New()
		defer mockDB.Close()

		err := RunETL(mockDB, "/fake/games.json", "valid/path", "valid/dir")

		assert.Error(t, err)
	})

	t.Run("Fails on step 2 persistMethods", func(t *testing.T) {
		mockDB, mockSQL, _ := sqlmock.New()
		defer mockDB.Close()

		gameContent := `[{"abbreviation": "R", "name": "Red", "generation": 1}]`
		gamesFile := testutil.CreateTempFile(t, "", "games_*.json", gameContent)
		defer gamesFile.Close()

		mockSQL.ExpectExec(`INSERT INTO .* \(abbreviation, name, generation\)`).
			WithArgs("R", "Red", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := RunETL(mockDB, gamesFile.Name(), "/fake/methods.json", "valid/dir")

		assert.Error(t, err)
		err = mockSQL.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Fails on step 3 persistNormalizedPokemon", func(t *testing.T) {
		mockDB, mockSQL, _ := sqlmock.New()
		defer mockDB.Close()

		gameContent := `[{"abbreviation": "R", "name": "Red", "generation": 1}]`
		gamesFile := testutil.CreateTempFile(t, "", "games_*.json", gameContent)
		defer gamesFile.Close()

		mockSQL.ExpectExec(`INSERT INTO .* \(abbreviation, name, generation\)`).
			WithArgs("R", "Red", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		methodContent := `{"WILD": "Wild encounter"}`
		methodsFile := testutil.CreateTempFile(t, "", "methods_*.json", methodContent)
		defer methodsFile.Close()

		mockSQL.ExpectExec(`INSERT INTO .* \(key, description\)`).
			WithArgs("WILD", "Wild encounter").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := RunETL(mockDB, gamesFile.Name(), methodsFile.Name(), "/fake/dir")

		assert.Error(t, err)
		err = mockSQL.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestPersistGames(t *testing.T) {
	t.Run("Fails on load execute", func(t *testing.T) {
		mockDB, _, _ := sqlmock.New()
		defer mockDB.Close()

		err := persistGames(mockDB, "/fake/path/games.json")

		assert.Error(t, err)
	})

	t.Run("Fails on save execute", func(t *testing.T) {
		mockDB, mockSQL, _ := sqlmock.New()
		defer mockDB.Close()

		gameContent := `[{"abbreviation": "R", "name": "Red", "generation": 1}]`
		file := testutil.CreateTempFile(t, "", "games_*.json", gameContent)
		defer file.Close()

		mockSQL.ExpectExec(`INSERT INTO .*`).
			WithArgs("R", "Red", 1).
			WillReturnError(errors.New("simulated database error"))

		err := persistGames(mockDB, file.Name())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "simulated database error")
	})
}

func TestPersistMethods(t *testing.T) {
	t.Run("Fails on load execute", func(t *testing.T) {
		mockDB, _, _ := sqlmock.New()
		defer mockDB.Close()

		err := persistMethods(mockDB, "/fake/path/methods.json")

		assert.Error(t, err)
	})

	t.Run("Fails on save execute", func(t *testing.T) {
		mockDB, mockSQL, _ := sqlmock.New()
		defer mockDB.Close()

		methodContent := `{"WILD": "Wild encounter"}`
		file := testutil.CreateTempFile(t, "", "methods_*.json", methodContent)
		defer file.Close()

		mockSQL.ExpectExec(`INSERT INTO .*`).
			WithArgs("WILD", "Wild encounter").
			WillReturnError(errors.New("database unavailable"))

		err := persistMethods(mockDB, file.Name())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database unavailable")
	})
}

func TestPersistNormalizedPokemon(t *testing.T) {
	t.Run("Fails on load execute", func(t *testing.T) {
		mockDB, _, _ := sqlmock.New()
		defer mockDB.Close()

		err := persistNormalizedPokemon(mockDB, "/fake/path/dir")

		assert.Error(t, err)
	})

	t.Run("Fails on save execute", func(t *testing.T) {
		mockDB, mockSQL, _ := sqlmock.New()
		defer mockDB.Close()

		pokemonDir := testutil.CreateTempDir(t, "", "pokemons_*")
		defer os.RemoveAll(pokemonDir)

		pokemonContent := `[{"number": 1, "name": "Bulbasaur", "form": "", "availability": [{"game": "R", "method": "WILD", "notes": ""}]}]`
		testutil.CreateFile(t, pokemonDir, "1.json", pokemonContent)

		mockSQL.ExpectQuery(`INSERT INTO .* \(number, name, form\)`).
			WithArgs(1, "Bulbasaur", "").
			WillReturnError(errors.New("constraint violation"))

		err := persistNormalizedPokemon(mockDB, pokemonDir)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "constraint violation")
	})
}
