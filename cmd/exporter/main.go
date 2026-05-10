package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/AdelinoLSN/pokemon-availability/internal/adapters/exporter"
	"github.com/AdelinoLSN/pokemon-availability/internal/adapters/repository"
	"github.com/AdelinoLSN/pokemon-availability/internal/adapters/source"
	"github.com/AdelinoLSN/pokemon-availability/internal/app"
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases"
)

func main() {
	log.Default().Println("Starting export...")

	app.LoadEnvironmentVariables()
	db, err := app.InitDatabaseConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	gamesFilepath := os.Getenv("APP_GAMES_JSON_FILEPATH")

	if err := runExporter(db, gamesFilepath); err != nil {
		panic(err)
	}

	log.Default().Println("Data exported successfully")
}

func runExporter(db *sql.DB, gamesFilepath string) error {
	games := loadGames(gamesFilepath)

	exportPokemonAvailabilityDetails(db, games)

	return nil
}

func loadGames(gamesFilepath string) []domain.Game {
	gameSource := source.NewJsonGameSource(gamesFilepath)

	loadGames := usecases.NewLoadGames(gameSource)

	games, err := loadGames.Execute()
	if err != nil {
		panic(err)
	}

	log.Default().Println("Reloaded games")

	return games
}

func exportPokemonAvailabilityDetails(db *sql.DB, games []domain.Game) {
	pokemonAvailabilityDetailRepository := repository.NewPostgresPokemonAvailabilityDetailRepository(db)
	csvPokemonAvailabilityDetailExporter := exporter.NewCsvPokemonAvailabilityDetailExporter()

	loadPokemonAvailabilityDetails := usecases.NewLoadPokemonAvailabilityDetails(pokemonAvailabilityDetailRepository)
	exportPokemonAvailabilityDetails := usecases.NewExportPokemonAvailabilityDetails(
		csvPokemonAvailabilityDetailExporter,
	)

	for i, game := range games {
		pokemonAvailabilityDetails, err := loadPokemonAvailabilityDetails.Execute(game.Abbreviation)
		if err != nil {
			panic(err)
		}

		err = exportPokemonAvailabilityDetails.Execute(i, game, pokemonAvailabilityDetails)
		if err != nil {
			panic(err)
		}

		log.Default().Printf("Exported data for game %s", game.Abbreviation)
	}
}
