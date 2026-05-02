package main

import (
	"log"

	"pokemon-availability/internal/adapters/database"
	"pokemon-availability/internal/adapters/file"
	"pokemon-availability/internal/usecase"
)

func main() {
	log.Println("Starting ETL process...")

	jsonReader := file.NewJSONReader()
	loader := usecase.NewLoader(jsonReader)

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.InitSchema(db); err != nil {
		log.Fatal(err)
	}

	methods, _ := loader.LoadMethods()
	games, _ := loader.LoadGames()
	pokemon, _ := loader.LoadPokemon()

	database.InsertMethods(db, methods)
	database.InsertGames(db, games)
	database.InsertPokemon(db, pokemon)

	log.Println("ETL finished with success")

  ExportToCSV(db, games)
}
