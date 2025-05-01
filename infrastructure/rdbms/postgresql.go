package rdbms

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Import driver PostgreSQL
)

var pgdatabase *sql.DB

// InitDB menginisialisasi koneksi database
func InitPG(host string, port int32, username string, password string, dbname string, sslmode string) error {
	var err error
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, username, password, dbname, sslmode,
	)

	pgdatabase, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// Cek koneksi database
	if err = pgdatabase.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	return nil
}

// GetDB mengembalikan instance database
func GetPG() *sql.DB {
	return pgdatabase
}
