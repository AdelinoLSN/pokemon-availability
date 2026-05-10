package source

import (
	"io/fs"
	"path/filepath"

	"github.com/AdelinoLSN/pokemon-availability/internal/dto"
	"github.com/AdelinoLSN/pokemon-availability/internal/infra/filesystem"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

var _ ports.PokemonSource = (*JsonPokemonSource)(nil)

type JsonPokemonSource struct {
	dirpath string
}

func NewJsonPokemonSource(dirpath string) *JsonPokemonSource {
	return &JsonPokemonSource{
		dirpath: dirpath,
	}
}

func (s *JsonPokemonSource) LoadPokemonsJson() ([]dto.PokemonJson, error) {
	var pokemonsJson []dto.PokemonJson

	err := filepath.WalkDir(
		s.dirpath,
		s.walkPokemonFiles(&pokemonsJson),
	)

	if err != nil {
		return nil, err
	}

	return pokemonsJson, nil
}

func (s *JsonPokemonSource) walkPokemonFiles(pokemonsJson *[]dto.PokemonJson) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !s.isJsonFile(d) {
			return nil
		}

		pokemon, err := s.readPokemonJsonFile(path)
		if err != nil {
			return err
		}

		*pokemonsJson = append(*pokemonsJson, pokemon...)

		return nil
	}
}

func (s *JsonPokemonSource) isJsonFile(d fs.DirEntry) bool {
	return !d.IsDir() && filepath.Ext(d.Name()) == ".json"
}

func (s *JsonPokemonSource) readPokemonJsonFile(path string) ([]dto.PokemonJson, error) {
	var pokemonsJson []dto.PokemonJson

	err := filesystem.ReadJson(path, &pokemonsJson)
	if err != nil {
		return nil, err
	}

	return pokemonsJson, nil
}
