package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTempDir(t *testing.T) {
	dir := CreateTempDir(t, "", "pokemon_*")
	defer os.RemoveAll(dir)

	info, err := os.Stat(dir)

	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestCreateTempFile(t *testing.T) {
	file := CreateTempFile(t, "", "pokemon_*.json", `{"name":"Pikachu"}`)
	defer os.Remove(file.Name())
	defer file.Close()

	data, err := os.ReadFile(file.Name())

	require.NoError(t, err)
	assert.Equal(t, `{"name":"Pikachu"}`, string(data))
}

func TestCreateFile(t *testing.T) {
	dir := t.TempDir()

	CreateFile(t, dir, "pokemon.json", `{"name":"Pikachu"}`)

	data, err := os.ReadFile(filepath.Join(dir, "pokemon.json"))
	require.NoError(t, err)
	assert.Equal(t, `{"name":"Pikachu"}`, string(data))
}
