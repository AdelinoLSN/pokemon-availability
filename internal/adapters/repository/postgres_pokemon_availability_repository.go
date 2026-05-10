package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
	"github.com/AdelinoLSN/pokemon-availability/internal/infra/database/postgres"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
)

var _ ports.PokemonAvailabilityRepository = (*PostgresPokemonAvailabilityRepository)(nil)

type PostgresPokemonAvailabilityRepository struct {
	db *sql.DB
}

func NewPostgresPokemonAvailabilityRepository(db *sql.DB) *PostgresPokemonAvailabilityRepository {
	return &PostgresPokemonAvailabilityRepository{
		db: db,
	}
}

func (r *PostgresPokemonAvailabilityRepository) SaveAll(pokemonAvailabilities []domain.PokemonAvailability) error {
	numRecords := len(pokemonAvailabilities)
	if numRecords == 0 {
		return nil
	}

	values := make([]string, 0, numRecords)
	args := make([]any, 0, numRecords*4)

	for i, availability := range pokemonAvailabilities {
		offset := i * 4

		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d)", offset+1, offset+2, offset+3, offset+4))

		args = append(args,
			availability.PokemonId,
			availability.GameAbbreviation,
			availability.MethodKey,
			availability.Note,
		)
	}

	query := fmt.Sprintf(
		`
		INSERT INTO %s (pokemon_id, game_abbreviation, method_key, note)
		VALUES %s
		ON CONFLICT (pokemon_id, game_abbreviation, method_key) DO NOTHING
		`,
		postgres.Tables.PokemonAvailability,
		strings.Join(values, ","),
	)

	_, err := r.db.Exec(query, args...)

	return err
}
