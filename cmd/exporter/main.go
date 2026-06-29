package main

import (
	"log"
	"os"

	"github.com/AdelinoLSN/pokemon-availability/internal/app"
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

	if err := app.RunExporter(db, gamesFilepath); err != nil {
		panic(err)
	}

	log.Default().Println("Data exported successfully")
}
