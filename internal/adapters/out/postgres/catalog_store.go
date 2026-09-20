package postgresadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/AdelinoLSN/pokemon-availability/internal/application/importcatalog"
	"github.com/AdelinoLSN/pokemon-availability/internal/domain"
)

var _ importcatalog.Store = (*CatalogStore)(nil)

type CatalogStore struct {
	db *sql.DB
}

func NewCatalogStore(db *sql.DB) *CatalogStore {
	return &CatalogStore{db: db}
}

func (s *CatalogStore) Upsert(ctx context.Context, catalog importcatalog.Catalog) error {
	transaction, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin catalog transaction: %w", err)
	}
	defer transaction.Rollback()

	if err := saveGames(ctx, transaction, catalog.Games); err != nil {
		return err
	}
	if err := saveMethods(ctx, transaction, catalog.Methods); err != nil {
		return err
	}
	for _, pokemon := range catalog.Pokemons {
		pokemonID, err := savePokemon(ctx, transaction, pokemon)
		if err != nil {
			return err
		}
		if err := saveAvailabilities(ctx, transaction, pokemonID, pokemon.Availabilities); err != nil {
			return err
		}
	}
	if _, err := transaction.ExecContext(ctx, "REFRESH MATERIALIZED VIEW mv_pokemon_availability_details"); err != nil {
		return fmt.Errorf("refresh availability projection: %w", err)
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit catalog transaction: %w", err)
	}
	return nil
}

func saveGames(ctx context.Context, transaction *sql.Tx, games []domain.Game) error {
	for _, game := range games {
		if _, err := transaction.ExecContext(ctx, `
			INSERT INTO games (abbreviation, name, generation)
			VALUES ($1, $2, $3)
			ON CONFLICT (abbreviation) DO UPDATE
			SET name = EXCLUDED.name, generation = EXCLUDED.generation
		`, game.Abbreviation, game.Name, game.Generation); err != nil {
			return fmt.Errorf("save game %s: %w", game.Abbreviation, err)
		}
	}
	return nil
}

func saveMethods(ctx context.Context, transaction *sql.Tx, methods []domain.Method) error {
	for _, method := range methods {
		if _, err := transaction.ExecContext(ctx, `
			INSERT INTO methods (key, description)
			VALUES ($1, $2)
			ON CONFLICT (key) DO UPDATE SET description = EXCLUDED.description
		`, method.Key, method.Description); err != nil {
			return fmt.Errorf("save method %s: %w", method.Key, err)
		}
	}
	return nil
}

func savePokemon(ctx context.Context, transaction *sql.Tx, pokemon importcatalog.Pokemon) (int, error) {
	var id int
	err := transaction.QueryRowContext(ctx, `
		INSERT INTO pokemon (number, name, form)
		VALUES ($1, $2, $3)
		ON CONFLICT (number, form) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, pokemon.Number, pokemon.Name, pokemon.Form).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("save pokemon %d/%s: %w", pokemon.Number, pokemon.Form, err)
	}
	return id, nil
}

func saveAvailabilities(
	ctx context.Context,
	transaction *sql.Tx,
	pokemonID int,
	availabilities []importcatalog.Availability,
) error {
	if len(availabilities) == 0 {
		return nil
	}
	values := make([]string, 0, len(availabilities))
	arguments := make([]any, 0, len(availabilities)*4)
	for index, availability := range availabilities {
		offset := index * 4
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d)", offset+1, offset+2, offset+3, offset+4))
		arguments = append(arguments, pokemonID, availability.GameAbbreviation, availability.MethodKey, availability.Note)
	}
	query := fmt.Sprintf(`
		INSERT INTO pokemon_availability (pokemon_id, game_abbreviation, method_key, note)
		VALUES %s
		ON CONFLICT (pokemon_id, game_abbreviation, method_key) DO UPDATE
		SET note = EXCLUDED.note
	`, strings.Join(values, ","))
	if _, err := transaction.ExecContext(ctx, query, arguments...); err != nil {
		return fmt.Errorf("save pokemon availabilities: %w", err)
	}
	return nil
}
