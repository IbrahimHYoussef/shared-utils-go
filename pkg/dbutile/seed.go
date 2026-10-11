package dbutile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	SEED_UP              = "-- +seed Up\n"
	SEED_DOWN            = "-- +seed Down\n"
	SEED_START_STATEMENT = "-- +seed StatementBegin\n"
	SEED_END_STATMENT    = "-- +seed StatementEnd\n"
)

// getStringsBetween returns the text between the first occurrence of start
// and the first occurrence of end that follows it. It returns "" when either
// marker is missing.
func getStringsBetween(txt, start, end string) string {

	if len(txt) == 0 {
		return ""
	}

	s := strings.Index(txt, start)
	if s < 0 {
		return ""
	}
	s += len(start)

	e := strings.Index(txt[s:], end)
	if e < 0 {
		return ""
	}

	return txt[s : s+e]
}

// MustSeedDatabase runs the "-- +seed Up" section of every .sql file in
// seedDirectory, each in its own transaction, and panics on any failure.
func MustSeedDatabase(ctx context.Context, pool *pgxpool.Pool, seedDirectory string) {
	mustRunSeedSection(ctx, pool, seedDirectory, SEED_UP)
}

// MustUnseedDatabase runs the "-- +seed Down" section of every .sql file in
// seedDirectory, each in its own transaction, and panics on any failure.
func MustUnseedDatabase(ctx context.Context, pool *pgxpool.Pool, seedDirectory string) {
	mustRunSeedSection(ctx, pool, seedDirectory, SEED_DOWN)
}

// mustRunSeedSection applies the statements between
// section+SEED_START_STATEMENT and SEED_END_STATMENT for every .sql file in
// seedDirectory, in file-name order. Files without that section are skipped.
func mustRunSeedSection(ctx context.Context, pool *pgxpool.Pool, seedDirectory, section string) {
	// make sure the directory exists
	if _, err := os.Stat(seedDirectory); err != nil {
		panic(fmt.Sprintf("Failed to find seed directory: %s", seedDirectory))
	}

	files, err := os.ReadDir(seedDirectory)
	if err != nil {
		panic(fmt.Sprintf("Failed to read seed directory %s: %v", seedDirectory, err))
	}

	// collect sql files in a stable order
	sqlFiles := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".sql" {
			continue
		}
		sqlFiles = append(sqlFiles, filepath.Join(seedDirectory, file.Name()))
	}
	sort.Strings(sqlFiles)

	for _, path := range sqlFiles {
		sqlByte, err := os.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("Failed to read seed file %s: %v", path, err))
		}

		statements := getStringsBetween(string(sqlByte), section+SEED_START_STATEMENT, SEED_END_STATMENT)
		if strings.TrimSpace(statements) == "" {
			continue
		}

		transaction, err := pool.Begin(ctx)
		if err != nil {
			panic(fmt.Sprintf("Failed to begin transaction for %s: %v", path, err))
		}

		if _, err := transaction.Exec(ctx, statements); err != nil {
			_ = transaction.Rollback(ctx)
			panic(fmt.Sprintf("Failed to run seed statements in %s: %v", path, err))
		}
		if err := transaction.Commit(ctx); err != nil {
			panic(fmt.Sprintf("Failed to commit seed statements in %s: %v", path, err))
		}
		logger.Info("Applied seed section", "file", path, "section", strings.TrimSpace(section))
	}
}
