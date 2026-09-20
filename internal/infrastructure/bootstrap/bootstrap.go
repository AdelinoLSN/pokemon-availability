package bootstrap

import (
	"context"
	"database/sql"

	"github.com/AdelinoLSN/pokemon-availability/internal/adapters/out/csv"
	"github.com/AdelinoLSN/pokemon-availability/internal/adapters/out/json"
	"github.com/AdelinoLSN/pokemon-availability/internal/adapters/out/postgres"
	"github.com/AdelinoLSN/pokemon-availability/internal/application/exportavailability"
	"github.com/AdelinoLSN/pokemon-availability/internal/application/importcatalog"
	"github.com/AdelinoLSN/pokemon-availability/internal/infrastructure/config"
	infrapostgres "github.com/AdelinoLSN/pokemon-availability/internal/infrastructure/postgres"
)

func OpenDatabase(ctx context.Context, configuration config.PostgresConfig) (*sql.DB, error) {
	return infrapostgres.NewConnection(ctx, configuration)
}

func ApplyMigrations(ctx context.Context, db *sql.DB) error {
	return infrapostgres.ApplyMigrations(ctx, db)
}

func NewImportCatalog(configuration config.Config, db *sql.DB) importcatalog.InputPort {
	source := jsoncatalog.NewCatalogSource(jsoncatalog.Paths{
		Games:    configuration.GamesJSONPath,
		Methods:  configuration.MethodsJSONPath,
		Pokemons: configuration.PokemonsJSONPath,
	})
	return importcatalog.New(source, postgresadapter.NewCatalogStore(db))
}

func NewExportAvailability(configuration config.Config, db *sql.DB) exportavailability.InputPort {
	return exportavailability.New(
		postgresadapter.NewAvailabilityReader(db),
		csvreport.NewAvailabilityWriter(configuration.OutputDirectory),
	)
}
