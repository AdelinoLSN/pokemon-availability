package database

import (
  "database/sql"
  "fmt"
  "strings"

  "pokemon-availability/internal/domain"
)

func InsertMethods(db *sql.DB, methods []domain.Method) error {
  const query = `
    INSERT INTO methods (key, description)
    VALUES ($1, $2)
    ON CONFLICT (key) DO NOTHING
  `

  for _, m := range methods {
    _, err := db.Exec(query, m.Key, m.Description)
    if err != nil {
      return fmt.Errorf("error inserting method %s: %w", m.Key, err)
    }
  }
  return nil
}

func InsertGames(db *sql.DB, games []domain.Game) error {
  const query = `
    INSERT INTO games (abbreviation, name, generation)
    VALUES ($1, $2, $3)
    ON CONFLICT (abbreviation) DO NOTHING
  `

  for _, g := range games {
    _, err := db.Exec(query, g.Abbreviation, g.Name, g.Generation)
    if err != nil {
      return fmt.Errorf("error inserting game %s: %w", g.Abbreviation, err)
    }
  }
  return nil
}

func InsertPokemon(db *sql.DB, pokemon []domain.Pokemon) error {
  tx, err := db.Begin()
  if err != nil {
    return fmt.Errorf("error starting transaction: %w", err)
  }
  defer tx.Rollback()

  const queryPokemon = `
    WITH inserted AS (
      INSERT INTO pokemon (number, name, form)
      VALUES ($1, $2, $3)
      ON CONFLICT (number, form) DO NOTHING
      RETURNING id
    )
    SELECT id FROM inserted
    UNION ALL
    SELECT id FROM pokemon WHERE number = $1 AND form = $3
    LIMIT 1;`

  const queryAvailability = `
    INSERT INTO pokemon_availability (pokemon_id, game_abbreviation, method_key, note)
    VALUES ($1, $2, $3, $4);`

  for _, p := range pokemon {
    var pokemonID int

    err := tx.QueryRow(queryPokemon, p.Number, p.Name, p.Form).Scan(&pokemonID)
    if err != nil {
      return fmt.Errorf("error while inserting or finding pokemon %s: %w", p.Name, err)
    }

    for _, a := range p.Availability {
      _, err := tx.Exec(queryAvailability, pokemonID, a.Game, a.Method, a.Notes)
      if err != nil {
        return fmt.Errorf("error while inserting availability for %s (Game: %s): %w", p.Name, a.Game, err)
      }
    }
  }

  if err := tx.Commit(); err != nil {
    return fmt.Errorf("error committing transaction: %w", err)
  }

  return nil
}

type PokemonCSVRow struct {
	Number  int
	Name    string
	Form    string
	Methods string
	Notes   string
}

func GetPokemonForGame(db *sql.DB, gameAbbreviation string) ([]PokemonCSVRow, error) {
	query := `
		SELECT
			number,
			name,
			COALESCE(form, '') AS form,
			STRING_AGG(method_key, '|') AS methods,
			STRING_AGG(COALESCE(note, ''), '|') AS notes
		FROM pokemon_full_view
		WHERE game_abbreviation = $1
		GROUP BY number, name, form
		ORDER BY number, form;
	`

	rows, err := db.Query(query, gameAbbreviation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []PokemonCSVRow
	for rows.Next() {
		var row PokemonCSVRow
		if err := rows.Scan(&row.Number, &row.Name, &row.Form, &row.Methods, &row.Notes); err != nil {
			return nil, err
		}

		row.Notes = strings.ReplaceAll(row.Notes, ",", ";")

		results = append(results, row)
	}

	return results, rows.Err()
}
