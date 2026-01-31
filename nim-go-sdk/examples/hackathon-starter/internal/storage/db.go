package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

// NewDB creates a new database connection and initializes the schema.
// It uses SQLite3 as the backing store.
func NewDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Use a context with timeout to validate the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	// Configure connection pool settings
	// These values are suitable for a small application; adjust based on workload
	db.SetMaxOpenConns(5)                  // Maximum number of open connections
	db.SetMaxIdleConns(5)                  // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Recycle connections every 5 minutes

	if err := initializeSchema(ctx, db); err != nil {
		db.Close()
		return nil, err
	}

	return &DB{conn: db}, nil
}

// initializeSchema creates the employees table if it doesn't exist
func initializeSchema(ctx context.Context, db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS employees (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		recipient TEXT UNIQUE NOT NULL,
		wage REAL NOT NULL DEFAULT 0,
		department TEXT NOT NULL DEFAULT ''
	);
	CREATE INDEX IF NOT EXISTS idx_employees_recipient ON employees(recipient);
	CREATE INDEX IF NOT EXISTS idx_employees_department ON employees(department);
	`

	if _, err := db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("error creating employees table: %w", err)
	}

	return nil
}

// Close closes the database connection
func (d *DB) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
