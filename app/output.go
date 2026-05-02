package main

import (
	"database/sql"
	"log"

	"pokemon-availability/internal/domain"
	"pokemon-availability/internal/adapters/file"
	"pokemon-availability/internal/usecase"
)

func ExportToCSV(db *sql.DB, games []domain.Game) {
	log.Println("Starting CSV Export process...")

	csvHandler := file.NewCSVHandler()

	exporter := usecase.NewExporter(db, csvHandler)

	if err := exporter.ExportGamesToCSV(games); err != nil {
		log.Fatalf("Error exporting CSVs: %v", err)
	}

	log.Println("CSV Export process finished with success")
}
