package main

import (
	"context"
	"log"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/importcatalog"
	"github.com/AdelinoLSN/pokemon-availability/internal/infrastructure/bootstrap"
	"github.com/AdelinoLSN/pokemon-availability/internal/infrastructure/config"
)

func main() {
	log.Default().Println("Starting ETL")

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
	if err := runEtl(ctx, bootstrap.NewImportCatalog(configuration, db)); err != nil {
		panic(err)
	}

	log.Default().Println("ETL finished")
}

func runEtl(ctx context.Context, useCase importcatalog.InputPort) error {
	return useCase.Execute(ctx, importcatalog.Command{})
}
