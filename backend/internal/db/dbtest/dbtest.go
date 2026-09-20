// Package dbtest starts an ephemeral PostgreSQL container for backend
// integration tests. It requires a working Docker daemon.
package dbtest

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/LinnikD/learn-the-grammar/backend/internal/db"
)

// StartPostgres starts an ephemeral PostgreSQL container, applies
// every embedded migration, and returns a connection pool to it. The
// container and pool are torn down via t.Cleanup.
func StartPostgres(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()

	container, err := postgres.Run(ctx, "postgres:17-alpine",
		postgres.WithDatabase("learn_the_grammar_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("starting postgres container: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Logf("terminating postgres container: %v", err)
		}
	})

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("reading postgres connection string: %v", err)
	}

	if err := db.MigrateUp(dsn); err != nil {
		t.Fatalf("applying migrations: %v", err)
	}

	pool, err := db.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("connecting to postgres: %v", err)
	}
	t.Cleanup(pool.Close)

	return pool
}
