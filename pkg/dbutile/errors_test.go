package dbutile

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestMapError(t *testing.T) {
	if err := MapError(nil); err != nil {
		t.Fatalf("nil: got %v, want nil", err)
	}
	if err := MapError(fmt.Errorf("wrapped: %w", pgx.ErrNoRows)); !errors.Is(err, ErrNotFound) {
		t.Fatalf("no rows: got %v, want ErrNotFound", err)
	}
	other := errors.New("boom")
	if err := MapError(other); err != other {
		t.Fatalf("other: got %v, want %v", err, other)
	}
}

func TestUniqueViolationConstraint(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23505", ConstraintName: "users_email_lower_key"}

	name, ok := UniqueViolationConstraint(fmt.Errorf("insert: %w", pgErr))
	if !ok || name != "users_email_lower_key" {
		t.Fatalf("unique violation: got %q %v", name, ok)
	}
	if !IsUniqueViolation(pgErr) {
		t.Fatalf("IsUniqueViolation: got false, want true")
	}
	if IsUniqueViolation(&pgconn.PgError{Code: "23503"}) {
		t.Fatalf("foreign key violation: got true, want false")
	}
	if IsUniqueViolation(errors.New("boom")) {
		t.Fatalf("plain error: got true, want false")
	}
}
