package main

import (
	"log"
	"os"

	"github.com/AdelinoLSN/pokemon-availability/internal/app"
)

func main() {
	log.Default().Println("Starting ETL")

	app.LoadEnvironmentVariables()
	db, err := app.InitDatabaseConnection()
	if err != nil {
		panic(err)
	}
	if err := app.InitDatabaseSchema(db); err != nil {
		panic(err)
	}
	defer db.Close()

	gamesFilepath := os.Getenv("APP_GAMES_JSON_FILEPATH")
	methodsFilepath := os.Getenv("APP_METHODS_JSON_FILEPATH")
	pokemonsJsonDirpath := os.Getenv("APP_POKEMONS_JSON_DIRPATH")

	if err := app.RunETL(db, gamesFilepath, methodsFilepath, pokemonsJsonDirpath); err != nil {
		panic(err)
	}

	log.Default().Println("ETL finished")
}
