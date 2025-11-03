package database

import (
    "context"
    "database/sql"
    "embed"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Migrate(ctx context.Context, db *sql.DB) error {
    b, err := migrationsFS.ReadFile("migrations/0001_init.sql")
    if err != nil {
        return err
    }
    _, err = db.ExecContext(ctx, string(b))
    return err
}

