package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresConnection() (*sql.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DB_DSN não definido")
	}

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return nil, err
		}

		if err = db.Ping(); err == nil {
			log.Println("Connected to database.")
			return db, nil
		}

		log.Println("Waiting database connection...")
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("não conseguiu conectar no banco")
}
