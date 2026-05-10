package app

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	log.Default().Println("Loaded environment variables")

	return nil
}
