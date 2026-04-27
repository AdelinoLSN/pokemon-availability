package usecase

import (
	"io/fs"
	"path/filepath"

	"pokemon-availability/internal/domain"
)

type JSONFileReader interface {
    Read(path string, v any) error
}

type Loader struct {
  jsonReader JSONFileReader
}

func NewLoader(j JSONFileReader) *Loader {
  return &Loader{j}
}

func (l *Loader) LoadMethods() ([]domain.Method, error) {
  var raw domain.MethodMap
  if err := l.jsonReader.Read("data/methods.json", &raw); err != nil {
    return nil, err
  }

  var methods []domain.Method
  for k, v := range raw {
    methods = append(methods, domain.Method{k, v})
  }

  return methods, nil
}

func (l *Loader) LoadGames() ([]domain.Game, error) {
  var games []domain.Game
  if err := l.jsonReader.Read("data/games.json", &games); err != nil {
    return nil, err
  }
  return games, nil
}

func (l *Loader) LoadPokemon() ([]domain.Pokemon, error) {
    var allPokemon []domain.Pokemon
    
    err := filepath.WalkDir("data/pokemon", func(path string, d fs.DirEntry, err error) error {
        if !d.IsDir() && filepath.Ext(d.Name()) == ".json" {
            var forms []domain.Pokemon
            if err := l.jsonReader.Read(path, &forms); err != nil {
                return err
            }
            allPokemon = append(allPokemon, forms...)
        }
        return nil
    })
    
    return allPokemon, err
}
