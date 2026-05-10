package postgres

import (
	"database/sql"
	"fmt"
)

var createTableGamesSql = fmt.Sprintf(
	`
	CREATE TABLE IF NOT EXISTS %s (
		abbreviation TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		generation INTEGER NOT NULL
	);
	`,
	Tables.Games,
)

var createTableMethodsSql = fmt.Sprintf(
	`
	CREATE TABLE IF NOT EXISTS %s (
		key TEXT PRIMARY KEY,
		description TEXT NOT NULL
	);
	`,
	Tables.Methods,
)

var createTablePokemonSql = fmt.Sprintf(
	`
	CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		number INTEGER NOT NULL,
		name TEXT NOT NULL,
		form TEXT,
		UNIQUE (number, form)
	);
	`,
	Tables.Pokemon,
)

var createTablePokemonAvailabilitySql = fmt.Sprintf(
	`
	CREATE TABLE IF NOT EXISTS %s (
		pokemon_id INTEGER NOT NULL,
		game_abbreviation TEXT NOT NULL,
		Method_key TEXT NOT NULL,
		note TEXT,
		UNIQUE ( pokemon_id, game_abbreviation, method_key ),
		FOREIGN KEY (pokemon_id) REFERENCES %s(id),
		FOREIGN KEY (game_abbreviation) REFERENCES %s(abbreviation),
		FOREIGN KEY (method_key) REFERENCES %s(key)
	);
	`,
	Tables.PokemonAvailability,
	Tables.Pokemon,
	Tables.Games,
	Tables.Methods,
)

var createViewPokemonAvailabilityDetailsSql = fmt.Sprintf(
	`
	CREATE OR REPLACE VIEW %s AS
	SELECT
		p.number AS number,
		p.name AS name,
		p.form,
		g.abbreviation AS game_abbreviation,
		g.name AS game,
		M.key AS method_key,
		M.description AS method_description,
		pa.note,
		p.id AS id
	FROM %s p
	JOIN %s pa ON pa.pokemon_id = p.id
	JOIN %s g ON g.abbreviation = pa.game_abbreviation
	JOIN %s m ON m.key = pa.method_key;
	`,
	Views.PokemonAvailabilityDetails,
	Tables.Pokemon,
	Tables.PokemonAvailability,
	Tables.Games,
	Tables.Methods,
)

var createViewPokemonAvailabilityDetailsInconsistencesSql = fmt.Sprintf(
	`
	CREATE OR REPLACE VIEW %s AS
	SELECT
		name,
		form,
		COUNT(DISTINCT number) AS total_numbers,
		ARRAY_AGG(DISTINCT number) AS numbers
	FROM %s
	GROUP BY name, form
	HAVING COUNT(DISTINCT number) > 1;
	`,
	Views.PokemonAvailabilityDetailsInconsistencies,
	Views.PokemonAvailabilityDetails,
)

var createMaterializedViewPokemonAvailabilityDetailsSql = fmt.Sprintf(
	`
	CREATE MATERIALIZED VIEW IF NOT EXISTS %s AS
	SELECT
		p.number AS number,
		p.name AS name,
		p.form,
		g.abbreviation AS game_abbreviation,
		g.name AS game,
		M.key AS method_key,
		M.description AS method_description,
		pa.note,
		p.id AS id
	FROM %s p
	JOIN %s pa ON pa.pokemon_id = p.id
	JOIN %s g ON g.abbreviation = pa.game_abbreviation
	JOIN %s m ON m.key = pa.method_key;
	`,
	MaterializedViews.PokemonAvailabilityDetails,
	Tables.Pokemon,
	Tables.PokemonAvailability,
	Tables.Games,
	Tables.Methods,
)

var createIndexesForMaterializedViewPokemonAvailabilityDetails = []string{
	fmt.Sprintf(
		`
		CREATE UNIQUE INDEX IF NOT EXISTS
		%s_unique_idx
		ON %s ( id, number, game_abbreviation, method_key);
		`,
		MaterializedViews.PokemonAvailabilityDetails,
		MaterializedViews.PokemonAvailabilityDetails,
	),

	fmt.Sprintf(
		`
		CREATE INDEX IF NOT EXISTS
		%s_number_idx
		ON %s (number);
		`,
		MaterializedViews.PokemonAvailabilityDetails,
		MaterializedViews.PokemonAvailabilityDetails,
	),

	fmt.Sprintf(
		`
		CREATE INDEX IF NOT EXISTS
		%s_name_idx
		ON %s (name);
		`,
		MaterializedViews.PokemonAvailabilityDetails,
		MaterializedViews.PokemonAvailabilityDetails,
	),

	fmt.Sprintf(
		`
		CREATE INDEX IF NOT EXISTS
		%s_game_idx
		ON %s (game_abbreviation);
		`,
		MaterializedViews.PokemonAvailabilityDetails,
		MaterializedViews.PokemonAvailabilityDetails,
	),
}

func InitSchema(db *sql.DB) error {
	queries := []string{
		// Tables
		createTableGamesSql,
		createTableMethodsSql,
		createTablePokemonSql,
		createTablePokemonAvailabilitySql,
		// Views
		createViewPokemonAvailabilityDetailsSql,
		createViewPokemonAvailabilityDetailsInconsistencesSql,
		// Materialized Views
		createMaterializedViewPokemonAvailabilityDetailsSql,
	}

	// Indexes
	queries = append(queries, createIndexesForMaterializedViewPokemonAvailabilityDetails...)

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}
