package main

import (
	"context"
	"log"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/exportavailability"
	"github.com/AdelinoLSN/pokemon-availability/internal/infrastructure/bootstrap"
	"github.com/AdelinoLSN/pokemon-availability/internal/infrastructure/config"
)

func main() {
	log.Default().Println("Starting export...")

	ctx := context.Background()
	configuration, err := config.Load()
	if err != nil {
		panic(err)
	}
	db, err := bootstrap.OpenDatabase(ctx, configuration.Postgres)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := bootstrap.ApplyMigrations(ctx, db); err != nil {
		panic(err)
	}

	if err := runExporter(ctx, bootstrap.NewExportAvailability(configuration, db)); err != nil {
		panic(err)
	}

	log.Default().Println("Data exported successfully")
}

func runExporter(ctx context.Context, useCase exportavailability.InputPort) error {
	return useCase.Execute(ctx, exportavailability.Command{})
}
