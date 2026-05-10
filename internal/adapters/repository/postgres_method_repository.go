package repository

import (
	"database/sql"
	"fmt"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/infra/database/postgres"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

var _ ports.MethodRepository = (*PostgresMethodRepository)(nil)

type PostgresMethodRepository struct {
	db *sql.DB
}

func NewPostgresMethodRepository(db *sql.DB) *PostgresMethodRepository {
	return &PostgresMethodRepository{
		db: db,
	}
}

func (r *PostgresMethodRepository) Save(method domain.Method) error {
	query := fmt.Sprintf(
		`
		INSERT INTO %s
		(key, description)
		VALUES ($1, $2)
		ON CONFLICT (key) DO NOTHING
		`,
		postgres.Tables.Methods,
	)

	_, err := r.db.Exec(query, method.Key, method.Description)

	return err
}

func (r *PostgresMethodRepository) SaveAll(methods []domain.Method) error {
	for _, method := range methods {
		err := r.Save(method)
		if err != nil {
			return err
		}
	}

	return nil
}
