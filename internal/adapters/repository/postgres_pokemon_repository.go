package repository

import (
	"database/sql"
	"fmt"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/infra/database/postgres"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

var _ ports.PokemonRepository = (*PostgresPokemonRepository)(nil)

type PostgresPokemonRepository struct {
	db *sql.DB
}

func NewPostgresPokemonRepository(db *sql.DB) *PostgresPokemonRepository {
	return &PostgresPokemonRepository{
		db: db,
	}
}

func (r *PostgresPokemonRepository) Save(pokemon domain.Pokemon) (int, error) {
	query := fmt.Sprintf(
		`
		INSERT INTO %s
		(number, name, form)
		VALUES ($1, $2, $3)
		ON CONFLICT (number, form) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
		`,
		postgres.Tables.Pokemon,
	)

	var id int

	err := r.db.QueryRow(
		query,
		pokemon.Number,
		pokemon.Name,
		pokemon.Form,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
