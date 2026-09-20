package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AdelinoLSN/pokemon-availability/internal/infrastructure/config"
	_ "github.com/lib/pq"
)

func NewConnection(ctx context.Context, config config.PostgresConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
