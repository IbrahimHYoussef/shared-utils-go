// Package testutils contains helpers for integration tests.
package testutils

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresTestCred contains credentials used to initialize a Postgres test
// container.
type PostgresTestCred struct {
	// DataBase is the database name created in the container.
	DataBase string
	// UserName is the database user name.
	UserName string
	// Password is the database user password.
	Password string
}

// PostgresContainer wraps a testcontainers Postgres container with its
// connection string.
type PostgresContainer struct {
	*postgres.PostgresContainer
	// ConnectionString is the PostgreSQL connection string with sslmode disabled.
	ConnectionString string
}

// CreatePGContainer starts a Postgres test container and returns it with a
// connection string.
//
// If image is empty, CreatePGContainer uses postgres:18-alpine. Callers are
// responsible for terminating the returned container when the test is done.
func CreatePGContainer(ctx context.Context, image string, creds *PostgresTestCred) (*PostgresContainer, error) {
	// get the env values for the test database

	// default container
	if len(image) == 0 {
		image = "postgres:18-alpine"
	}
	// default creds
	if creds != nil {
		creds = &PostgresTestCred{
			DataBase: "",
			UserName: "",
			Password: "",
		}
	}

	pgContainer, err := postgres.Run(
		ctx,
		image,
		postgres.WithDatabase(creds.DataBase),
		postgres.WithUsername(creds.UserName),
		postgres.WithPassword(creds.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		return nil, err
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}
