package dbutile

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
