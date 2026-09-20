CREATE TABLE IF NOT EXISTS games (
    abbreviation TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    generation INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS methods (
    key TEXT PRIMARY KEY,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS pokemon (
    id SERIAL PRIMARY KEY,
    number INTEGER NOT NULL,
    name TEXT NOT NULL,
    form TEXT,
    UNIQUE (number, form)
);

CREATE TABLE IF NOT EXISTS pokemon_availability (
    pokemon_id INTEGER NOT NULL,
    game_abbreviation TEXT NOT NULL,
    method_key TEXT NOT NULL,
    note TEXT,
    UNIQUE (pokemon_id, game_abbreviation, method_key),
    FOREIGN KEY (pokemon_id) REFERENCES pokemon(id),
    FOREIGN KEY (game_abbreviation) REFERENCES games(abbreviation),
    FOREIGN KEY (method_key) REFERENCES methods(key)
);

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_pokemon_availability_details AS
SELECT
    p.number,
    p.name,
    p.form,
    g.abbreviation AS game_abbreviation,
    m.key AS method_key,
    pa.note
FROM pokemon p
JOIN pokemon_availability pa ON pa.pokemon_id = p.id
JOIN games g ON g.abbreviation = pa.game_abbreviation
JOIN methods m ON m.key = pa.method_key;

CREATE UNIQUE INDEX IF NOT EXISTS mv_pokemon_availability_details_unique_idx
ON mv_pokemon_availability_details (number, form, game_abbreviation, method_key);

CREATE INDEX IF NOT EXISTS mv_pokemon_availability_details_game_idx
ON mv_pokemon_availability_details (game_abbreviation);
