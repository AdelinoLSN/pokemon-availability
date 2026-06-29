package app

import (
	"database/sql"
	"log"

	"github.com/AdelinoLSN/pokemon-availability/internal/adapters/exporter"
	"github.com/AdelinoLSN/pokemon-availability/internal/adapters/repository"
	"github.com/AdelinoLSN/pokemon-availability/internal/adapters/source"
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases"
)

func RunExporter(db *sql.DB, gamesFilepath string) error {
	games, err := loadGames(gamesFilepath)
	if err != nil {
		return err
	}

	if err := exportPokemonAvailabilityDetails(db, games); err != nil {
		return err
	}

	return nil
}

func loadGames(gamesFilepath string) ([]domain.Game, error) {
	gameSource := source.NewJsonGameSource(gamesFilepath)

	loadGames := usecases.NewLoadGames(gameSource)

	games, err := loadGames.Execute()
	if err != nil {
		return nil, err
	}

	log.Default().Println("Reloaded games")

	return games, nil
}

func exportPokemonAvailabilityDetails(db *sql.DB, games []domain.Game) error {
	pokemonAvailabilityDetailRepository := repository.NewPostgresPokemonAvailabilityDetailRepository(db)
	csvPokemonAvailabilityDetailExporter := exporter.NewCsvPokemonAvailabilityDetailExporter()

	loadPokemonAvailabilityDetails := usecases.NewLoadPokemonAvailabilityDetails(pokemonAvailabilityDetailRepository)
	exportPokemonAvailabilityDetails := usecases.NewExportPokemonAvailabilityDetails(
		csvPokemonAvailabilityDetailExporter,
	)

	for i, game := range games {
		pokemonAvailabilityDetails, err := loadPokemonAvailabilityDetails.Execute(game.Abbreviation)
		if err != nil {
			return err
		}

		err = exportPokemonAvailabilityDetails.Execute(i, game, pokemonAvailabilityDetails)
		if err != nil {
			return err
		}

		log.Default().Printf("Exported data for game %s", game.Abbreviation)
	}

	return nil
}
