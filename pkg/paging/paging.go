// Package paging contains helpers for cursor-paged list endpoints: clamping
// the limit query parameter, encoding and decoding opaque keyset cursors, and
// building the page response.
package paging

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// ErrInvalidCursor is returned by DecodeCursor when the cursor was not
// produced by EncodeCursor for the same key type.
var ErrInvalidCursor = errors.New("invalid cursor")

// Page is the response body of a list endpoint.
type Page[T any] struct {
	// Items is the rows of this page. It is encoded as [] when empty.
	Items []T `json:"items"`
	// NextCursor is the cursor of the next page, nil on the last page.
	NextCursor *string `json:"next_cursor"`
}

// Cursor is a keyset cursor for lists ordered by one sort key and then by id.
//
// Lists ordered by more than one key can pass their own struct to EncodeCursor
// and DecodeCursor instead.
type Cursor[K any] struct {
	// Key is the sort key of the last row of the page.
	Key K `json:"k"`
	// ID is the id of the last row of the page.
	ID uuid.UUID `json:"id"`
}

// ClampLimit returns the page size for the raw limit query value.
//
// The limit is never rejected: an absent value, 0, or a value that is not a
// number returns def, any other value below 1 returns 1, and a value above max
// returns max. Fractions are rounded down first.
func ClampLimit(raw string, def, max int) int {
	n, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil || math.IsNaN(n) {
		return def
	}
	n = math.Floor(n)
	switch {
	case n == 0:
		return def
	case n < 1:
		return 1
	case n > float64(max):
		return max
	}
	return int(n)
}

// EncodeCursor returns key as an opaque cursor string.
func EncodeCursor[K any](key K) (string, error) {
	data, err := json.Marshal(key)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

// DecodeCursor returns the key stored in cursor.
//
// An empty cursor means the first page and returns nil with no error. A cursor
// that cannot be decoded into K returns ErrInvalidCursor.
func DecodeCursor[K any](cursor string) (*K, error) {
	if cursor == "" {
		return nil, nil
	}
	data, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, ErrInvalidCursor
	}
	key := new(K)
	if err := json.Unmarshal(data, key); err != nil {
		return nil, ErrInvalidCursor
	}
	return key, nil
}

// NewPage builds a Page from rows fetched with a limit of limit+1.
//
// When rows holds more than limit items there is a next page: the extra row is
// dropped and NextCursor is built by passing the last kept row to keyOf.
func NewPage[T any, K any](rows []T, limit int, keyOf func(T) K) (Page[T], error) {
	if len(rows) <= limit {
		if rows == nil {
			rows = []T{}
		}
		return Page[T]{Items: rows}, nil
	}

	rows = rows[:limit]
	cursor, err := EncodeCursor(keyOf(rows[limit-1]))
	if err != nil {
		return Page[T]{}, err
	}
	return Page[T]{Items: rows, NextCursor: &cursor}, nil
}
