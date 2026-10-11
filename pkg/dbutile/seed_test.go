package dbutile

import (
	"context"
	"slices"
	"testing"

	"github.com/IbrahimHYoussef/shared-utils-go/pkg/testutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
)

type StringBetweenCases struct {
	Input    string
	Expected string
	Start    string
	End      string
}

func Test_getStringsBetween(t *testing.T) {

	cases := map[string]StringBetweenCases{
		"empty string":                         {Input: "", Expected: "", Start: "Start", End: "End"},
		"normal Strings no \n":                 {Input: "START THIS IS what i need  END", Expected: " THIS IS what i need  ", Start: "START", End: "END"},
		"strings with \n and \t":               {Input: "START THIS IS\n what i \t need  END", Expected: " THIS IS\n what i \t need  ", Start: "START", End: "END"},
		"multiline sperators strings with \n":  {Input: "START\nBEGIN\n THIS IS\n what i \t need  \nEND", Expected: " THIS IS\n what i \t need  \n", Start: "START\nBEGIN\n", End: "END"},
		"end marker also appears before start": {Input: "START A END\nDOWN B END", Expected: " B ", Start: "DOWN", End: "END"},
	}

	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("getStringsBetween panicked: %v", r)
				}
			}()
			txt := getStringsBetween(testCase.Input, testCase.Start, testCase.End)

			if txt != testCase.Expected {
				t.Errorf("got |%s| want |%s|", txt, testCase.Expected)
			}
		})
	}
}

const seedTestDir = "../../test_data/seed_test"

// newSeedTestPool starts a Postgres container with an empty test_table and
// returns a pool connected to it. Both are cleaned up when the test ends.
func newSeedTestPool(t *testing.T) (context.Context, *pgxpool.Pool) {
	t.Helper()
	testcontainers.SkipIfProviderIsNotHealthy(t)
	ctx := context.Background()

	pgContainer, err := testutils.CreatePGContainer(ctx, "", nil)
	if err != nil {
		t.Fatalf("Failed To Start Container: %v", err)
	}
	t.Cleanup(func() { _ = pgContainer.Terminate(ctx) })

	pool, err := pgxpool.New(ctx, pgContainer.ConnectionString)
	if err != nil {
		t.Fatalf("Failed To Connect To Database: %v", err)
	}
	t.Cleanup(pool.Close)

	if _, err := pool.Exec(ctx, "CREATE TABLE test_table (id INT PRIMARY KEY, name TEXT NOT NULL)"); err != nil {
		t.Fatalf("Failed To Create Table: %v", err)
	}
	return ctx, pool
}

// seededNames returns the names in test_table ordered by id.
func seededNames(t *testing.T, ctx context.Context, pool *pgxpool.Pool) []string {
	t.Helper()
	rows, err := pool.Query(ctx, "SELECT name FROM test_table ORDER BY id")
	if err != nil {
		t.Fatalf("Failed To Query test_table: %v", err)
	}
	defer rows.Close()
	names := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("Failed To Scan: %v", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows: %v", err)
	}
	return names
}

func Test_MustSeedDatabase(t *testing.T) {
	t.Run("seeds the Up section of every sql file", func(t *testing.T) {
		ctx, pool := newSeedTestPool(t)

		MustSeedDatabase(ctx, pool, seedTestDir)

		got := seededNames(t, ctx, pool)
		want := []string{"alice", "bob", "carol"}
		if !slices.Equal(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})

	t.Run("panics when the seed directory does not exist", func(t *testing.T) {
		ctx, pool := newSeedTestPool(t)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for missing seed directory")
			}
		}()
		MustSeedDatabase(ctx, pool, "../../test_data/does_not_exist")
	})

	t.Run("panics when a seed statement fails", func(t *testing.T) {
		ctx, pool := newSeedTestPool(t)
		if _, err := pool.Exec(ctx, "DROP TABLE test_table"); err != nil {
			t.Fatalf("drop: %v", err)
		}
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic when the insert fails")
			}
		}()
		MustSeedDatabase(ctx, pool, seedTestDir)
	})
}

func Test_MustUnseedDatabase(t *testing.T) {
	t.Run("removes the seeded rows using the Down section", func(t *testing.T) {
		ctx, pool := newSeedTestPool(t)
		MustSeedDatabase(ctx, pool, seedTestDir)
		if len(seededNames(t, ctx, pool)) != 3 {
			t.Fatal("precondition: seed did not insert rows")
		}

		MustUnseedDatabase(ctx, pool, seedTestDir)

		if got := seededNames(t, ctx, pool); len(got) != 0 {
			t.Fatalf("expected empty table after unseed, got %v", got)
		}
	})

	t.Run("panics when the seed directory does not exist", func(t *testing.T) {
		ctx, pool := newSeedTestPool(t)
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for missing seed directory")
			}
		}()
		MustUnseedDatabase(ctx, pool, "../../test_data/does_not_exist")
	})
}
