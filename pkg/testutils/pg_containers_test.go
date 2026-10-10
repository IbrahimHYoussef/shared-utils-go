package testutils_test

import (
	"context"
	"testing"

	"github.com/IbrahimHYoussef/shared-utils-go/pkg/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
)

func TestCreatePGContainer(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	ctx := context.Background()

	cases := map[string]*testutils.PostgresTestCred{
		"nil creds use defaults":         nil,
		"given creds are kept":           {DataBase: "app", UserName: "app_user", Password: "secret"},
		"empty fields fall back to test": {DataBase: "app"},
	}
	for name, creds := range cases {
		t.Run(name, func(t *testing.T) {
			container, err := testutils.CreatePGContainer(ctx, "", creds)
			if err != nil {
				t.Fatalf("CreatePGContainer: %v", err)
			}
			t.Cleanup(func() { _ = container.Terminate(ctx) })

			conn, err := pgx.Connect(ctx, container.ConnectionString)
			if err != nil {
				t.Fatalf("connect: %v", err)
			}
			defer conn.Close(ctx)

			var database, user, version string
			if err := conn.QueryRow(ctx, "select current_database(), current_user, current_setting('server_version_num')").Scan(&database, &user, &version); err != nil {
				t.Fatalf("query: %v", err)
			}
			wantDB, wantUser := "test", "test"
			if creds != nil && creds.DataBase != "" {
				wantDB = creds.DataBase
			}
			if creds != nil && creds.UserName != "" {
				wantUser = creds.UserName
			}
			if database != wantDB || user != wantUser {
				t.Fatalf("connected to %s as %s, want %s as %s", database, user, wantDB, wantUser)
			}
			if version < "180000" {
				t.Fatalf("server_version_num %s, want Postgres 18 by default", version)
			}
		})
	}
}
