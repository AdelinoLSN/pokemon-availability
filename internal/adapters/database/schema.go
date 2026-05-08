package database

import "database/sql"

func InitSchema(db *sql.DB) error {
	queries := []string{

		`CREATE TABLE IF NOT EXISTS methods (
			key TEXT PRIMARY KEY,
			description TEXT NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS games (
			abbreviation TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			generation INTEGER NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS pokemon (
			id SERIAL PRIMARY KEY,
			number INTEGER NOT NULL,
			name TEXT NOT NULL,
			form TEXT,
			UNIQUE (number, form)
		);`,

		`CREATE TABLE IF NOT EXISTS pokemon_availability (
			pokemon_id INTEGER NOT NULL,
			game_abbreviation TEXT NOT NULL,
			method_key TEXT NOT NULL,
			note TEXT,
      UNIQUE (pokemon_id, game_abbreviation, method_key),
			FOREIGN KEY (pokemon_id) REFERENCES pokemon(id),
			FOREIGN KEY (game_abbreviation) REFERENCES games(abbreviation),
			FOREIGN KEY (method_key) REFERENCES methods(key)
		);`,

		`CREATE OR REPLACE VIEW pokemon_full_view AS
		SELECT
			p.number as number,
			p.name AS name,
			p.form,
			g.abbreviation AS game_abbreviation,
			g.name AS game,
			m.key as method_key,
			m.description AS method_description,
      pa.note,
      p.id as id
		FROM pokemon p
		JOIN pokemon_availability pa ON pa.pokemon_id = p.id
		JOIN games g ON g.abbreviation = pa.game_abbreviation
		JOIN methods m ON m.key = pa.method_key;`,

		`CREATE OR REPLACE VIEW public.pokemon_number_inconsistencies AS
		SELECT
				name,
				form,
				COUNT(DISTINCT number) AS qtd_numeros,
				ARRAY_AGG(DISTINCT number) AS numeros
		FROM public.pokemon_full_view
		GROUP BY name, form
		HAVING COUNT(DISTINCT number) > 1;`,

		`CREATE MATERIALIZED VIEW IF NOT EXISTS pokemon_full_mv AS
		SELECT
		p.number as number,
		p.name AS name,
		p.form,
		g.abbreviation AS game_abbreviation,
		g.name AS game,
		m.key as method_key,
		m.description AS method_description,
		pa.note,
		p.id as id
		FROM pokemon p
		JOIN pokemon_availability pa ON pa.pokemon_id = p.id
		JOIN games g ON g.abbreviation = pa.game_abbreviation
		JOIN methods m ON m.key = pa.method_key;`,

		`CREATE UNIQUE INDEX IF NOT EXISTS pokemon_full_mv_unique_idx
		ON pokemon_full_mv (
		id,
		number,
		game_abbreviation,
		method_key
		);`,

		`CREATE INDEX IF NOT EXISTS pokemon_full_mv_game_idx
		ON pokemon_full_mv (game_abbreviation);`,

		`CREATE INDEX IF NOT EXISTS pokemon_full_mv_number_idx
		ON pokemon_full_mv (number);`,

		`CREATE INDEX IF NOT EXISTS pokemon_full_mv_name_idx
		ON pokemon_full_mv (name);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}
