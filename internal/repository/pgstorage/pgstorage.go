package pgstorage

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

const (
	ErrIsExistCode = "23505"

	migrationsTable = "schema_migrations"
	schemaName      = "public"
	migrationsPath  = "./migrations"
)

// PGStorage wraps sql.DB with PostgreSQL connection pool and migration management.
type PGStorage struct {
	db *sql.DB
}

// New creates PostgreSQL storage with automatic schema migrations.
// 1. Establishes pgx connection pool
// 2. Initializes migration driver with schema_migrations table
// 3. Applies all pending migrations from ./migrations directory
// 4. Returns sql.DB compatible storage instance
func New(connStr string) (*PGStorage, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDBFromPool(pool)

	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: migrationsTable,
		SchemaName:      schemaName,
	})
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+absPath, "postgres", driver)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return &PGStorage{
		db: db,
	}, nil
}

// Shutdown closes database connection pool gracefully.
// Call during application shutdown to release resources.
func (s *PGStorage) Shutdown() error {
	return s.db.Close()
}
