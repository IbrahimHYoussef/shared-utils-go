// Package dbutile contains PostgreSQL connection, migration, transaction, and
// logging helpers built on pgx.
package dbutile

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectToDatabase creates a pgx connection pool for dsn, verifies it with
// Ping, and returns the ready pool.
//
// ConnectToDatabase configures MinConns to 1 and MaxConns to 10. It logs and
// panics if the DSN cannot be parsed, the pool cannot be created, or the ping
// fails.
func ConnectToDatabase(ctx context.Context, dsn string) *pgxpool.Pool {
	// create pgx connection pool with config
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error("Failed to parse database config", "dsn", dsn, "error", err)
		panic(err)
	}

	// MinConns ensures at least 1 connection is established immediately
	poolConfig.MinConns = 1
	poolConfig.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Error("Failed to connect to database", "dsn", dsn, "error", err)
		panic(err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		logger.Error("Failed To Connect To database", "dsn", dsn, "error", err)
		panic(err)
	}

	logger.Info("Successfully connected to database")
	return pool
}
