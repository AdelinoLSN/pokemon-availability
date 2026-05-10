package repository

import (
	"database/sql"
	"fmt"

	"github.com/AdelinoLSN/pokemon-availability/internal/infra/database/postgres"
	"github.com/AdelinoLSN/pokemon-availability/internal/ports"
	"github.com/AdelinoLSN/pokemon-availability/internal/usecases/models"
)

var _ ports.PokemonAvailabilityDetailRepository = (*PostgresPokemonAvailabilityDetailRepository)(nil)

type PostgresPokemonAvailabilityDetailRepository struct {
	db *sql.DB
}

func NewPostgresPokemonAvailabilityDetailRepository(db *sql.DB) *PostgresPokemonAvailabilityDetailRepository {
	return &PostgresPokemonAvailabilityDetailRepository{
		db: db,
	}
}

func (r *PostgresPokemonAvailabilityDetailRepository) RefreshMaterializedView() error {
	query := fmt.Sprintf(
		`REFRESH MATERIALIZED VIEW %s`,
		postgres.MaterializedViews.PokemonAvailabilityDetails,
	)

	_, err := r.db.Exec(query)

	return err
}

func (r *PostgresPokemonAvailabilityDetailRepository) LoadByGameAbbreviation(
	gameAbbreviation string,
) ([]models.PokemonAvailabilityDetail, error) {
	query := fmt.Sprintf(
		`
		SELECT
			number,
			name,
			form,
			game_abbreviation,
			game,
			method_key,
			method_description,
			note,
			id
		FROM %s
		WHERE game_abbreviation = $1
		`,
		postgres.MaterializedViews.PokemonAvailabilityDetails,
	)

	rows, err := r.db.Query(
		query,
		gameAbbreviation,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pokemonAvailabilitiesDetails []models.PokemonAvailabilityDetail

	for rows.Next() {
		var pokemonAvailabilityDetails models.PokemonAvailabilityDetail

		err := rows.Scan(
			&pokemonAvailabilityDetails.Number,
			&pokemonAvailabilityDetails.Name,
			&pokemonAvailabilityDetails.Form,
			&pokemonAvailabilityDetails.GameAbbreviation,
			&pokemonAvailabilityDetails.Game,
			&pokemonAvailabilityDetails.MethodKey,
			&pokemonAvailabilityDetails.MethodDescription,
			&pokemonAvailabilityDetails.Note,
			&pokemonAvailabilityDetails.Id,
		)

		if err != nil {
			return nil, err
		}

		pokemonAvailabilitiesDetails = append(
			pokemonAvailabilitiesDetails,
			pokemonAvailabilityDetails,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pokemonAvailabilitiesDetails, nil
}
