package connections

import (
	"database/sql"
	"log"
	"os"
)

func InitPostgres() (*sql.DB, error) {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@postgres:5432/hotty_delivery?sslmode=disable"
	}
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
