package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" database/sql driver used by goose
	"github.com/pressly/goose/v3"
)

// MigrationsDir is the absolute path to the goose SQL migrations,
// resolved relative to this source file so callers do not depend on
// the process's working directory.
var MigrationsDir = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "db", "migrations")
}()

// MigrateUp applies every pending migration in dir to the database at dsn.
func MigrateUp(dsn, dir string) error {
	return withGooseDB(dsn, func(sqlDB *sql.DB) error {
		return goose.Up(sqlDB, dir)
	})
}

// MigrateDown rolls back the most recently applied migration in dir.
func MigrateDown(dsn, dir string) error {
	return withGooseDB(dsn, func(sqlDB *sql.DB) error {
		return goose.Down(sqlDB, dir)
	})
}

// MigrateStatus writes the status of every migration in dir to goose's
// configured logger (standard output by default).
func MigrateStatus(dsn, dir string) error {
	return withGooseDB(dsn, func(sqlDB *sql.DB) error {
		return goose.Status(sqlDB, dir)
	})
}

// withGooseDB opens a database/sql connection scoped to PostgreSQL only:
// unlike the goose CLI, this backend does not need goose's support for
// every dialect it knows how to migrate.
func withGooseDB(dsn string, fn func(*sql.DB) error) error {
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
