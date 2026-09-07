package dbutile

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// MustMigrateUp runs SQL migrations from migrationsDir and panics on error.
//
// MustMigrateUp uses goose with the postgres dialect and the provided pgx pool.
func MustMigrateUp(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) {
	logger.Info("Running database migrations", "dir", migrationsDir)

	// Convert pgxpool to database/sql for goose
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	// Set goose dialect
	if err := goose.SetDialect("postgres"); err != nil {
		panic(fmt.Sprintf("failed to set goose dialect: %v", err))
	}
	path, err := os.Getwd()
	if err != nil {
		panic("Failed to get working dir ")
	}
	logger.Info("Path Check", "path", path)
	// Run migrations
	if err := goose.UpContext(ctx, db, migrationsDir); err != nil {
		// panic(fmt.Sprintf("failed to run migrations: %v", err))
		formatted := err.Error()
		formatted = strings.ReplaceAll(formatted, "\\n", "\n")
		formatted = strings.ReplaceAll(formatted, "\\t", "\t")

		panic("failed to run migraiton:\n" + formatted)
	}

	logger.Info("Successfully completed all migrations")
}
