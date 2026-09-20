package jsoncatalog

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCatalogSourceTranslatesJSONAtTheBoundary(t *testing.T) {
	root := t.TempDir()
	pokemonPath := filepath.Join(root, "pokemon")
	if err := os.Mkdir(pokemonPath, 0o755); err != nil {
		t.Fatal(err)
	}
	gamesPath := filepath.Join(root, "games.json")
	methodsPath := filepath.Join(root, "methods.json")
	writeFile(t, gamesPath, `[{"abbreviation":"R","name":"Red","generation":1}]`)
	writeFile(t, methodsPath, `{"WILD":"Wild encounter"}`)
	writeFile(t, filepath.Join(pokemonPath, "0001.json"), `[{"number":1,"name":"Bulbasaur","form":"","availability":[{"game":"R","method":"WILD","notes":""}]}]`)

	catalog, err := NewCatalogSource(Paths{Games: gamesPath, Methods: methodsPath, Pokemons: pokemonPath}).Load(context.Background())

	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if catalog.Games[0].Name != "Red" || catalog.Methods[0].Key != "WILD" || catalog.Pokemons[0].Availabilities[0].GameAbbreviation != "R" {
		t.Fatalf("catalog = %#v", catalog)
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
