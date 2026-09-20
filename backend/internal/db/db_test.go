package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LinnikD/learn-the-grammar/backend/internal/db"
	"github.com/LinnikD/learn-the-grammar/backend/internal/db/dbtest"
)

// TestPing proves the full pipeline end to end: an ephemeral
// PostgreSQL container is started, goose migrations are applied
// against it, and a sqlc-generated query runs successfully over a
// pgx connection pool.
func TestPing(t *testing.T) {
	pool := dbtest.StartPostgres(t)

	ok, err := db.New(pool).Ping(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int32(1), ok)
}
