package source

import (
	"io/fs"
	"path/filepath"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/dto"
	"github.com/AdelinoLSN/pokemon-availability/internal/infra/filesystem"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
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

func (s *JsonPokemonSource) LoadNormalizedPokemons() ([]models.NormalizedPokemon, error) {
	var normalizedPokemons []models.NormalizedPokemon

	err := filepath.WalkDir(
		s.dirpath,
		s.walkPokemonFiles(&normalizedPokemons),
	)

	if err != nil {
		return nil, err
	}

	return normalizedPokemons, nil
}

func (s *JsonPokemonSource) walkPokemonFiles(normalizedPokemons *[]models.NormalizedPokemon) fs.WalkDirFunc {
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

		for _, pokemonJson := range pokemon {
			*normalizedPokemons = append(
				*normalizedPokemons,
				s.buildNormalizedPokemon(pokemonJson),
			)
		}

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

func (s *JsonPokemonSource) buildNormalizedPokemon(pokemonJson dto.PokemonJson) models.NormalizedPokemon {
	availabilities := make([]domain.PokemonAvailability, 0, len(pokemonJson.Availability))

	for _, availabilityJson := range pokemonJson.Availability {
		availabilities = append(availabilities, domain.PokemonAvailability{
			GameAbbreviation: availabilityJson.Game,
			MethodKey:        availabilityJson.Method,
			Note:             availabilityJson.Notes,
		})
	}

	return models.NormalizedPokemon{
		Pokemon: domain.Pokemon{
			Number: pokemonJson.Number,
			Name:   pokemonJson.Name,
			Form:   pokemonJson.Form,
		},
		Availabilities: availabilities,
	}
}
