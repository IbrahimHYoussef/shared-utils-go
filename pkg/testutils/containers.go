package testutils

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresTestCred struct {
	DataBase string
	UserName string
	Password string
}

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

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
