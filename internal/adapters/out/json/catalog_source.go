package jsoncatalog

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/importcatalog"
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/infrastructure/filesystem"
)

var _ importcatalog.Source = (*CatalogSource)(nil)

type Paths struct {
	Games    string
	Methods  string
	Pokemons string
}

type CatalogSource struct {
	paths Paths
}

type gameJSON struct {
	Abbreviation string `json:"abbreviation"`
	Name         string `json:"name"`
	Generation   int    `json:"generation"`
}

type pokemonJSON struct {
	Number       int                `json:"number"`
	Name         string             `json:"name"`
	Form         string             `json:"form"`
	Availability []availabilityJSON `json:"availability"`
}

type availabilityJSON struct {
	Game   string `json:"game"`
	Method string `json:"method"`
	Notes  string `json:"notes"`
}

func NewCatalogSource(paths Paths) *CatalogSource {
	return &CatalogSource{paths: paths}
}

func (s *CatalogSource) Load(ctx context.Context) (importcatalog.Catalog, error) {
	if err := ctx.Err(); err != nil {
		return importcatalog.Catalog{}, err
	}
	games, err := s.loadGames()
	if err != nil {
		return importcatalog.Catalog{}, err
	}
	if err := ctx.Err(); err != nil {
		return importcatalog.Catalog{}, err
	}
	methods, err := s.loadMethods()
	if err != nil {
		return importcatalog.Catalog{}, err
	}
	pokemons, err := s.loadPokemons(ctx)
	if err != nil {
		return importcatalog.Catalog{}, err
	}

	return importcatalog.Catalog{
		Games:    games,
		Methods:  methods,
		Pokemons: pokemons,
	}, nil
}

func (s *CatalogSource) loadGames() ([]domain.Game, error) {
	var source []gameJSON
	if err := filesystem.ReadJSON(s.paths.Games, &source); err != nil {
		return nil, fmt.Errorf("read games: %w", err)
	}

	games := make([]domain.Game, 0, len(source))
	for _, game := range source {
		games = append(games, domain.Game{
			Abbreviation: game.Abbreviation,
			Name:         game.Name,
			Generation:   game.Generation,
		})
	}
	return games, nil
}

func (s *CatalogSource) loadMethods() ([]domain.Method, error) {
	var source map[string]string
	if err := filesystem.ReadJSON(s.paths.Methods, &source); err != nil {
		return nil, fmt.Errorf("read methods: %w", err)
	}

	keys := make([]string, 0, len(source))
	for key := range source {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	methods := make([]domain.Method, 0, len(source))
	for _, key := range keys {
		methods = append(methods, domain.Method{Key: key, Description: source[key]})
	}
	return methods, nil
}

func (s *CatalogSource) loadPokemons(ctx context.Context) ([]importcatalog.Pokemon, error) {
	pokemons := make([]importcatalog.Pokemon, 0)
	err := filepath.WalkDir(s.paths.Pokemons, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			return nil
		}

		var source []pokemonJSON
		if err := filesystem.ReadJSON(path, &source); err != nil {
			return fmt.Errorf("read pokemon file %s: %w", path, err)
		}
		for _, pokemon := range source {
			availabilities := make([]importcatalog.Availability, 0, len(pokemon.Availability))
			for _, availability := range pokemon.Availability {
				availabilities = append(availabilities, importcatalog.Availability{
					GameAbbreviation: availability.Game,
					MethodKey:        availability.Method,
					Note:             availability.Notes,
				})
			}
			pokemons = append(pokemons, importcatalog.Pokemon{
				Number:         pokemon.Number,
				Name:           pokemon.Name,
				Form:           pokemon.Form,
				Availabilities: availabilities,
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read pokemons: %w", err)
	}
	return pokemons, nil
}
