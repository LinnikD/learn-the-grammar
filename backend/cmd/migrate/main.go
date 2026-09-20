// Command migrate applies, rolls back, or reports the status of the
// PostgreSQL schema migrations embedded from backend/internal/db/migrations.
//
// It intentionally wraps only the goose library functions this backend
// needs instead of depending on goose's own CLI (cmd/goose), which
// links in a driver for every dialect goose supports.
package main

import (
	"fmt"
	"os"

	"github.com/LinnikD/learn-the-grammar/backend/internal/db"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: migrate <up|down|status>")
	}

	dsn := os.Getenv("LTG_DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("LTG_DATABASE_URL is not set")
	}

	switch args[0] {
	case "up":
		return db.MigrateUp(dsn)
	case "down":
		return db.MigrateDown(dsn)
	case "status":
		return db.MigrateStatus(dsn)
	default:
		return fmt.Errorf("unknown command %q: usage: migrate <up|down|status>", args[0])
	}
}
