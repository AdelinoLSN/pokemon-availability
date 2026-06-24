package app

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() error {
	if err := godotenv.Load(); err != nil {
		log.Default().Println("Failed to load .env file: using system environment variables")
		return err
	}

	log.Default().Println("Loaded environment variables from .env file")

	return nil
}
