package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUsers, downCreateUsers)
}

func upCreateUsers(ctx context.Context, tx *sql.Tx) error {
	// Create the users table
	_, err := tx.ExecContext(ctx, `
	CREATE TABLE users (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		username TEXT NOT NULL UNIQUE,
		emailAddress TEXT NOT NULL UNIQUE
	);
	`)
	return err
}

func downCreateUsers(ctx context.Context, tx *sql.Tx) error {
	// Drop the users table
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS users;`)
	return err
}
