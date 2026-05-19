package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/FranciscoHonorat/books-api/internal/infra/sqlc"
	_ "github.com/lib/pq"
)

// Database encapsulates the database connection.
type Database struct {
	db *sql.DB
}

// NewConnection creates a new database connection.
func NewConnection(dsn string) (*Database, error) {
	//Open a new database connection using the provided DSN.
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	//Pool settings. (MaxOpenConns, MaxIdleConns, ConnMaxLifetime)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	//Ping the database to ensure the connection is established.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	//Return the database instance and nil error.
	return &Database{db: db}, nil
}

// HealthCheck checks the database connection health.
func (d *Database) HealthCheck(ctx context.Context) error {
	//db.PingContext checks the database connection health using the provided context.
	if err := d.db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}

// Close closes the database connection.
func (d *Database) Close() error {
	//Close the database connection and return any error that occurs.
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}

// GetDB returns the underlying sql.DB instance for executing queries.
func (d *Database) GetDB() *sql.DB {
	return d.db
}

// Queries returns a new instance of sqlc.Queries for executing type-safe queries.
func (d *Database) Queries() *sqlc.Queries {
	return sqlc.New(d.db)
}
