package dbutile

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// uniqueViolationCode is the PostgreSQL SQLSTATE for unique_violation.
const uniqueViolationCode = "23505"

// ErrNotFound is returned by MapError when a query returned no rows.
var ErrNotFound = errors.New("not found")

// MapError converts pgx.ErrNoRows to ErrNotFound and returns every other error
// unchanged, including nil.
func MapError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

// IsUniqueViolation reports whether err is a PostgreSQL unique_violation.
func IsUniqueViolation(err error) bool {
	_, ok := UniqueViolationConstraint(err)
	return ok
}

// UniqueViolationConstraint returns the name of the violated constraint or
// index when err is a PostgreSQL unique_violation.
//
// UniqueViolationConstraint returns an empty string and false for any other
// error.
func UniqueViolationConstraint(err error) (string, bool) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != uniqueViolationCode {
		return "", false
	}
	return pgErr.ConstraintName, true
}
