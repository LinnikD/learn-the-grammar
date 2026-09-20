package db

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" database/sql driver used by goose
	"github.com/pressly/goose/v3"
)

// migrationsFS embeds the SQL migrations into the binary at build
// time. A path resolved at runtime (e.g. via runtime.Caller) would
// break once the binary is built with -trimpath or run from outside
// its original build directory.
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

// migrationsDir is the directory within migrationsFS that goose reads
// migrations from.
const migrationsDir = "migrations"

// MigrateUp applies every pending migration to the database at dsn.
func MigrateUp(dsn string) error {
	return withGooseDB(dsn, func(sqlDB *sql.DB) error {
		return goose.Up(sqlDB, migrationsDir)
	})
}

// MigrateDown rolls back the most recently applied migration.
func MigrateDown(dsn string) error {
	return withGooseDB(dsn, func(sqlDB *sql.DB) error {
		return goose.Down(sqlDB, migrationsDir)
	})
}

// MigrateStatus writes the status of every migration to goose's
// configured logger (standard output by default).
func MigrateStatus(dsn string) error {
	return withGooseDB(dsn, func(sqlDB *sql.DB) error {
		return goose.Status(sqlDB, migrationsDir)
	})
}

// withGooseDB opens a database/sql connection scoped to PostgreSQL only:
// unlike the goose CLI, this backend does not need goose's support for
// every dialect it knows how to migrate.
func withGooseDB(dsn string, fn func(*sql.DB) error) error {
	goose.SetBaseFS(migrationsFS)
	defer goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("setting goose dialect: %w", err)
	}

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("opening migration connection: %w", err)
	}
	defer func() { _ = sqlDB.Close() }()

	return fn(sqlDB)
}
