package app

import (
	"database/sql"
	"log"

	"github.com/AdelinoLSN/pokemon-availability/internal/adapters/repository"
	"github.com/AdelinoLSN/pokemon-availability/internal/adapters/source"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases"
)

func RunETL(db *sql.DB, gamesFilepath string, methodsFilepath string, pokemonsJsonDirpath string) error {
	if err := persistGames(db, gamesFilepath); err != nil {
		return err
	}

	if err := persistMethods(db, methodsFilepath); err != nil {
		return err
	}

	if err := persistNormalizedPokemon(db, pokemonsJsonDirpath); err != nil {
		return err
	}

	return nil
}

func persistGames(db *sql.DB, gamesFilepath string) error {
	gameSource := source.NewJsonGameSource(gamesFilepath)
	gameRepository := repository.NewPostgresGameRepository(db)

	loadGames := usecases.NewLoadGames(gameSource)
	saveGames := usecases.NewSaveGames(gameRepository)

	games, err := loadGames.Execute()
	if err != nil {
		return err
	}

	if err := saveGames.Execute(games); err != nil {
		return err
	}

	log.Default().Println("Saved games in database")

	return nil
}

func persistMethods(db *sql.DB, methodsFilepath string) error {
	methodSource := source.NewJsonMethodSource(methodsFilepath)
	methodRepository := repository.NewPostgresMethodRepository(db)

	loadMethods := usecases.NewLoadMethods(methodSource)
	saveMethods := usecases.NewSaveMethods(methodRepository)

	methods, err := loadMethods.Execute()
	if err != nil {
		return err
	}

	if err := saveMethods.Execute(methods); err != nil {
		return err
	}

	log.Default().Println("Saved methods in database")

	return nil
}

func persistNormalizedPokemon(db *sql.DB, pokemonsJsonDirpath string) error {
	pokemonSource := source.NewJsonPokemonSource(pokemonsJsonDirpath)
	pokemonRepository := repository.NewPostgresPokemonRepository(db)
	pokemonAvailabilityRepository := repository.NewPostgresPokemonAvailabilityRepository(db)
	pokemonAvailabilityDetailRepository := repository.NewPostgresPokemonAvailabilityDetailRepository(db)

	loadNormalizedPokemons := usecases.NewNormalizePokemon(pokemonSource)
	saveNormalizedPokemon := usecases.NewSaveNormalizedPokemons(
		pokemonRepository,
		pokemonAvailabilityRepository,
		pokemonAvailabilityDetailRepository,
	)

	normalizedPokemons, err := loadNormalizedPokemons.Execute()
	if err != nil {
		return err
	}

	if err := saveNormalizedPokemon.Execute(normalizedPokemons); err != nil {
		return err
	}

	log.Default().Println("Persisted normalized Pokémon")

	return nil
}
