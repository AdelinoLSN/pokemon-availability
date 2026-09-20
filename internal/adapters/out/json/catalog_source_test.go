package jsoncatalog

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCatalogSourceMapsJSONIntoTheApplicationCatalog(t *testing.T) {
	root := t.TempDir()
	games := filepath.Join(root, "games.json")
	methods := filepath.Join(root, "methods.json")
	pokemons := filepath.Join(root, "pokemon")
	if err := os.Mkdir(pokemons, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, games, `[{"abbreviation":"R","name":"Red","generation":1}]`)
	writeFile(t, methods, `{"WILD":"Wild encounter"}`)
	writeFile(t, filepath.Join(pokemons, "0001.json"), `[{"number":1,"name":"Bulbasaur","form":"","availability":[{"game":"R","method":"WILD","notes":""}]}]`)

	catalog, err := NewCatalogSource(Paths{
		Games: games, Methods: methods, Pokemons: pokemons,
	}).Load(context.Background())

	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if catalog.Games[0].Name != "Red" || catalog.Methods[0].Key != "WILD" {
		t.Fatalf("catalog = %#v", catalog)
	}
	pokemon := catalog.Pokemons[0]
	if pokemon.Number != 1 || pokemon.Availabilities[0].GameAbbreviation != "R" {
		t.Fatalf("pokemon = %#v", pokemon)
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
