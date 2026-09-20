package postgresadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/exportavailability"
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
)

var _ exportavailability.Reader = (*AvailabilityReader)(nil)

type AvailabilityReader struct {
	db *sql.DB
}

func NewAvailabilityReader(db *sql.DB) *AvailabilityReader {
	return &AvailabilityReader{db: db}
}

func (r *AvailabilityReader) ListGames(ctx context.Context) ([]domain.Game, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT abbreviation, name, generation
		FROM games
		ORDER BY generation, abbreviation
	`)
	if err != nil {
		return nil, fmt.Errorf("list games: %w", err)
	}
	defer rows.Close()
	games := make([]domain.Game, 0)
	for rows.Next() {
		var game domain.Game
		if err := rows.Scan(&game.Abbreviation, &game.Name, &game.Generation); err != nil {
			return nil, fmt.Errorf("scan game: %w", err)
		}
		games = append(games, game)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate games: %w", err)
	}
	return games, nil
}

func (r *AvailabilityReader) ByGame(ctx context.Context, gameAbbreviation string) ([]exportavailability.AvailabilityRow, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT number, name, form, method_key, note
		FROM mv_pokemon_availability_details
		WHERE game_abbreviation = $1
		ORDER BY number, form, method_key
	`, gameAbbreviation)
	if err != nil {
		return nil, fmt.Errorf("load availability for %s: %w", gameAbbreviation, err)
	}
	defer rows.Close()
	availabilities := make([]exportavailability.AvailabilityRow, 0)
	for rows.Next() {
		var row exportavailability.AvailabilityRow
		if err := rows.Scan(&row.Number, &row.Name, &row.Form, &row.MethodKey, &row.Note); err != nil {
			return nil, fmt.Errorf("scan availability: %w", err)
		}
		availabilities = append(availabilities, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate availability: %w", err)
	}
	return availabilities, nil
}
