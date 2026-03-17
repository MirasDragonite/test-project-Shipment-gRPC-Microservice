package postgres

import (
	"database/sql"
	"fmt"
	"os"
)

func NewDB() (*sql.DB, error) {

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")
	sslMode := os.Getenv("POSTGRES_SSLMODE")

	// Default values if we didn't get any
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if sslMode == "" {
		sslMode = "disable"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	// Checking DB
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}
