package repository

import (
	"database/sql"
	"fmt"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/infra/database/postgres"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

var _ ports.GameRepository = (*PostgresGameRepository)(nil)

type PostgresGameRepository struct {
	db *sql.DB
}

func NewPostgresGameRepository(db *sql.DB) *PostgresGameRepository {
	return &PostgresGameRepository{
		db: db,
	}
}

func (r *PostgresGameRepository) GetAll() ([]domain.Game, error) {
	query := fmt.Sprintf(
		`
		SELECT abbreviation, name, generation
		FROM %s
		ORDER BY generation, abbreviation
		`,
		postgres.Tables.Games,
	)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []domain.Game

	for rows.Next() {
		var game domain.Game

		err := rows.Scan(
			&game.Abbreviation,
			&game.Name,
			&game.Generation,
		)

		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

func (r *PostgresGameRepository) Save(game domain.Game) error {
	query := fmt.Sprintf(
		`
		INSERT INTO %s
		(abbreviation, name, generation)
		VALUES ($1, $2, $3)
		ON CONFLICT (abbreviation) DO NOTHING
		`,
		postgres.Tables.Games,
	)

	_, err := r.db.Exec(query, game.Abbreviation, game.Name, game.Generation)

	return err
}

func (r *PostgresGameRepository) SaveAll(games []domain.Game) error {
	for _, game := range games {
		err := r.Save(game)
		if err != nil {
			return err
		}
	}

	return nil
}
