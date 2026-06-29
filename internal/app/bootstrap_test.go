package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadEnvironmentVariables(t *testing.T) {
	workingDir, err := os.Getwd()
	require.NoError(t, err)

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte("POKEMON_TEST_ENV=loaded\n"), 0644))
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(workingDir)
	require.NoError(t, os.Unsetenv("POKEMON_TEST_ENV"))
	defer os.Unsetenv("POKEMON_TEST_ENV")

	err = LoadEnvironmentVariables()

	assert.NoError(t, err)
	assert.Equal(t, "loaded", os.Getenv("POKEMON_TEST_ENV"))
}

func TestLoadEnvironmentVariables_ReturnsErrorWhenDotEnvIsMissing(t *testing.T) {
	workingDir, err := os.Getwd()
	require.NoError(t, err)

	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(workingDir)

	err = LoadEnvironmentVariables()

	assert.Error(t, err)
}

func TestInitDatabaseConnection_ReturnsConnectionError(t *testing.T) {
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "1")
	t.Setenv("DB_USERNAME", "user")
	t.Setenv("DB_PASSWORD", "password")
	t.Setenv("DB_NAME", "pokemon")

	db, err := InitDatabaseConnection()

	assert.Nil(t, db)
	assert.Error(t, err)
}

func TestInitDatabaseSchema(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	for i := 0; i < 11; i++ {
		mock.ExpectExec(`(?s).*`).WillReturnResult(sqlmock.NewResult(0, 0))
	}

	err = InitDatabaseSchema(db)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
