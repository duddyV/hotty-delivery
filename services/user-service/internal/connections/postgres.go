package connections

import (
	"database/sql"
	"log"
	"os"
)

func InitPostgres() (*sql.DB, error) {
	dsn := os.Getenv("DSN")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Soccessfully connected to PostgreSQL")
	return db, nil
}
