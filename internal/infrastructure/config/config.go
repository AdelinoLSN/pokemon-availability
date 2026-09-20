package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GamesJSONPath    string
	MethodsJSONPath  string
	PokemonsJSONPath string
	OutputDirectory  string
	Postgres         PostgresConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load environment file: %w", err)
	}

	config := Config{
		GamesJSONPath:    os.Getenv("APP_GAMES_JSON_FILEPATH"),
		MethodsJSONPath:  os.Getenv("APP_METHODS_JSON_FILEPATH"),
		PokemonsJSONPath: os.Getenv("APP_POKEMONS_JSON_DIRPATH"),
		OutputDirectory:  os.Getenv("APP_OUTPUT_DIRECTORY"),
		Postgres: PostgresConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_NAME"),
		},
	}
	if config.OutputDirectory == "" {
		config.OutputDirectory = ".outputs"
	}
	return config, config.validate()
}

func (c Config) validate() error {
	required := map[string]string{
		"APP_GAMES_JSON_FILEPATH":   c.GamesJSONPath,
		"APP_METHODS_JSON_FILEPATH": c.MethodsJSONPath,
		"APP_POKEMONS_JSON_DIRPATH": c.PokemonsJSONPath,
		"DB_HOST":                   c.Postgres.Host,
		"DB_PORT":                   c.Postgres.Port,
		"DB_USERNAME":               c.Postgres.Username,
		"DB_PASSWORD":               c.Postgres.Password,
		"DB_NAME":                   c.Postgres.Database,
	}
	for name, value := range required {
		if value == "" {
			return fmt.Errorf("required environment variable %s is not set", name)
		}
	}
	return nil
}
