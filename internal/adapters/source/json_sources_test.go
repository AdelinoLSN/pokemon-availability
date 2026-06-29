package source

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonGameSource_LoadGames(t *testing.T) {
	path := filepath.Join(t.TempDir(), "games.json")
	require.NoError(t, os.WriteFile(path, []byte(`[{"abbreviation":"R","name":"Red","generation":1}]`), 0644))

	games, err := NewJsonGameSource(path).LoadGames()

	require.NoError(t, err)
	assert.Equal(t, []domain.Game{{Abbreviation: "R", Name: "Red", Generation: 1}}, games)
}

func TestJsonGameSource_LoadGames_ReturnsReadError(t *testing.T) {
	games, err := NewJsonGameSource(filepath.Join(t.TempDir(), "missing.json")).LoadGames()

	assert.Nil(t, games)
	assert.Error(t, err)
}

func TestJsonMethodSource_LoadMethods(t *testing.T) {
	path := filepath.Join(t.TempDir(), "methods.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"WILD":"Wild encounter"}`), 0644))

	methods, err := NewJsonMethodSource(path).LoadMethods()

	require.NoError(t, err)
	assert.Equal(t, []domain.Method{{Key: "WILD", Description: "Wild encounter"}}, methods)
}

func TestJsonMethodSource_LoadMethods_ReturnsReadError(t *testing.T) {
	methods, err := NewJsonMethodSource(filepath.Join(t.TempDir(), "missing.json")).LoadMethods()

	assert.Nil(t, methods)
	assert.Error(t, err)
}

func TestJsonPokemonSource_LoadNormalizedPokemons(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "ignore.txt"), []byte("ignored"), 0644))
	require.NoError(t, os.WriteFile(
		filepath.Join(dir, "pokemon.json"),
		[]byte(`[{"number":25,"name":"Pikachu","form":"","availability":[{"game":"R","method":"STARTER","notes":"Only one"}]}]`),
		0644,
	))

	pokemons, err := NewJsonPokemonSource(dir).LoadNormalizedPokemons()

	require.NoError(t, err)
	require.Len(t, pokemons, 1)
	assert.Equal(t, domain.Pokemon{Number: 25, Name: "Pikachu", Form: ""}, pokemons[0].Pokemon)
	assert.Equal(t, []domain.PokemonAvailability{
		{GameAbbreviation: "R", MethodKey: "STARTER", Note: "Only one"},
	}, pokemons[0].Availabilities)
}

func TestJsonPokemonSource_LoadNormalizedPokemons_ReturnsInvalidJsonError(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "pokemon.json"), []byte(`not-json`), 0644))

	pokemons, err := NewJsonPokemonSource(dir).LoadNormalizedPokemons()

	assert.Nil(t, pokemons)
	assert.Error(t, err)
}

func TestJsonPokemonSource_LoadNormalizedPokemons_ReturnsWalkError(t *testing.T) {
	pokemons, err := NewJsonPokemonSource(filepath.Join(t.TempDir(), "missing")).LoadNormalizedPokemons()

	assert.Nil(t, pokemons)
	assert.Error(t, err)
}
