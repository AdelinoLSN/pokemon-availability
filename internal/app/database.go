package app

import (
	"database/sql"
	"log"

	"github.com/AdelinoLSN/pokemon-availability/internal/infra/database/postgres"
)

func InitDatabaseConnection() (*sql.DB, error) {
	db, err := postgres.NewPostgresConnection()
	if err != nil {
		return nil, err
	}

	log.Default().Println("Initialized database")

	return db, err
}

func InitDatabaseSchema(db *sql.DB) error {
	if err := postgres.InitSchema(db); err != nil {
		return err
	}

	log.Default().Println("Initialized database schema")

	return nil
}
