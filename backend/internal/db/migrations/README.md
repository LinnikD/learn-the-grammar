# Database migrations

Goose SQL migrations for the application database. There is no product
schema here yet — FR issues that need persistent storage add their own
migrations to this directory.

## Adding a migration

Create a new `NNNNN_description.sql` file (goose also accepts
`go run github.com/pressly/goose/v3/cmd/goose create <name> sql` from
this directory to generate the boilerplate) with `-- +goose Up` and
`-- +goose Down` sections, for example:

```sql
-- +goose Up
CREATE TABLE example (
    id uuid PRIMARY KEY
);

-- +goose Down
DROP TABLE example;
```

## Running migrations

Use the Makefile targets, which read the database DSN from
`LTG_DATABASE_URL`:

```sh
make db-migrate-up
make db-migrate-down
make db-migrate-status
```

Migrations are applied with `backend/cmd/migrate`, a thin wrapper
around the `goose` library scoped to PostgreSQL only (the `goose` CLI
pulls in drivers for every dialect it supports, which this project
does not need).

`*.sql` files in this directory are embedded into the `migrate`
binary and the test suite at build time (see `//go:embed` in
`backend/internal/db/migrate.go`) instead of being read from disk at
run time, so migrations keep working regardless of the process's
working directory or how the binary was built (including
`-trimpath`). Rebuild after adding or editing a migration for the
change to take effect.
