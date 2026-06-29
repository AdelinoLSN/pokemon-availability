package usecases_test

import (
	"errors"
	"testing"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
	"github.com/stretchr/testify/assert"
)

type gameSourceStub struct {
	games []domain.Game
	err   error
}

func (s gameSourceStub) LoadGames() ([]domain.Game, error) {
	return s.games, s.err
}

type methodSourceStub struct {
	methods []domain.Method
	err     error
}

func (s methodSourceStub) LoadMethods() ([]domain.Method, error) {
	return s.methods, s.err
}

type pokemonSourceStub struct {
	pokemons []models.NormalizedPokemon
	err      error
}

func (s pokemonSourceStub) LoadNormalizedPokemons() ([]models.NormalizedPokemon, error) {
	return s.pokemons, s.err
}

type gameRepositoryStub struct {
	saved []domain.Game
	err   error
}

func (r *gameRepositoryStub) GetAll() ([]domain.Game, error) {
	return nil, nil
}

func (r *gameRepositoryStub) Save(game domain.Game) error {
	r.saved = append(r.saved, game)
	return r.err
}

func (r *gameRepositoryStub) SaveAll(games []domain.Game) error {
	r.saved = games
	return r.err
}

type methodRepositoryStub struct {
	saved []domain.Method
	err   error
}

func (r *methodRepositoryStub) Save(method domain.Method) error {
	r.saved = append(r.saved, method)
	return r.err
}

func (r *methodRepositoryStub) SaveAll(methods []domain.Method) error {
	r.saved = methods
	return r.err
}

type pokemonAvailabilityDetailRepositoryStub struct {
	details []models.PokemonAvailabilityDetail
	err     error
}

func (r *pokemonAvailabilityDetailRepositoryStub) RefreshMaterializedView() error {
	return nil
}

func (r *pokemonAvailabilityDetailRepositoryStub) LoadByGameAbbreviation(
	string,
) ([]models.PokemonAvailabilityDetail, error) {
	return r.details, r.err
}

type pokemonAvailabilityDetailExporterStub struct {
	path    string
	details []models.PokemonAvailabilityDetail
	err     error
}

func (e *pokemonAvailabilityDetailExporterStub) ExportPokemonAvailabilityDetails(
	path string,
	details []models.PokemonAvailabilityDetail,
) error {
	e.path = path
	e.details = details
	return e.err
}

func TestLoadGames(t *testing.T) {
	expected := []domain.Game{{Abbreviation: "R", Name: "Red", Generation: 1}}

	games, err := usecases.NewLoadGames(gameSourceStub{games: expected}).Execute()

	assert.NoError(t, err)
	assert.Equal(t, expected, games)
}

func TestLoadGames_ReturnsSourceError(t *testing.T) {
	expectedErr := errors.New("source failed")

	games, err := usecases.NewLoadGames(gameSourceStub{err: expectedErr}).Execute()

	assert.Nil(t, games)
	assert.ErrorIs(t, err, expectedErr)
}

func TestLoadMethods(t *testing.T) {
	expected := []domain.Method{{Key: "WILD", Description: "Wild encounter"}}

	methods, err := usecases.NewLoadMethods(methodSourceStub{methods: expected}).Execute()

	assert.NoError(t, err)
	assert.Equal(t, expected, methods)
}

func TestNormalizePokemon(t *testing.T) {
	expected := []models.NormalizedPokemon{
		{
			Pokemon: domain.Pokemon{Number: 25, Name: "Pikachu"},
		},
	}

	pokemons, err := usecases.NewNormalizePokemon(pokemonSourceStub{pokemons: expected}).Execute()

	assert.NoError(t, err)
	assert.Equal(t, expected, pokemons)
}

func TestSaveGames(t *testing.T) {
	games := []domain.Game{{Abbreviation: "R", Name: "Red", Generation: 1}}
	repository := &gameRepositoryStub{}

	err := usecases.NewSaveGames(repository).Execute(games)

	assert.NoError(t, err)
	assert.Equal(t, games, repository.saved)
}

func TestSaveGames_ReturnsRepositoryError(t *testing.T) {
	expectedErr := errors.New("save failed")
	repository := &gameRepositoryStub{err: expectedErr}

	err := usecases.NewSaveGames(repository).Execute([]domain.Game{{Abbreviation: "R"}})

	assert.ErrorIs(t, err, expectedErr)
}

func TestSaveMethods(t *testing.T) {
	methods := []domain.Method{{Key: "WILD", Description: "Wild encounter"}}
	repository := &methodRepositoryStub{}

	err := usecases.NewSaveMethods(repository).Execute(methods)

	assert.NoError(t, err)
	assert.Equal(t, methods, repository.saved)
}

func TestLoadPokemonAvailabilityDetails(t *testing.T) {
	expected := []models.PokemonAvailabilityDetail{{Number: 25, Name: "Pikachu"}}
	repository := &pokemonAvailabilityDetailRepositoryStub{details: expected}

	details, err := usecases.NewLoadPokemonAvailabilityDetails(repository).Execute("R")

	assert.NoError(t, err)
	assert.Equal(t, expected, details)
}

func TestLoadPokemonAvailabilityDetails_ReturnsRepositoryError(t *testing.T) {
	expectedErr := errors.New("load failed")
	repository := &pokemonAvailabilityDetailRepositoryStub{err: expectedErr}

	details, err := usecases.NewLoadPokemonAvailabilityDetails(repository).Execute("R")

	assert.Nil(t, details)
	assert.ErrorIs(t, err, expectedErr)
}

func TestExportPokemonAvailabilityDetails(t *testing.T) {
	exporter := &pokemonAvailabilityDetailExporterStub{}
	details := []models.PokemonAvailabilityDetail{{Number: 25, Name: "Pikachu"}}
	game := domain.Game{Abbreviation: "R"}

	err := usecases.NewExportPokemonAvailabilityDetails(exporter).Execute(7, game, details)

	assert.NoError(t, err)
	assert.Equal(t, ".outputs/007_R.csv", exporter.path)
	assert.Equal(t, details, exporter.details)
}

func TestExportPokemonAvailabilityDetails_ReturnsExporterError(t *testing.T) {
	expectedErr := errors.New("export failed")
	exporter := &pokemonAvailabilityDetailExporterStub{err: expectedErr}

	err := usecases.NewExportPokemonAvailabilityDetails(exporter).Execute(0, domain.Game{Abbreviation: "R"}, nil)

	assert.ErrorIs(t, err, expectedErr)
}
