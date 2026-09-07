// Package mapping contains conversion helpers between standard Go types and
// pgx pgtype values.
package mapping

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// PgUUIDFromGoogle converts a google/uuid UUID to a valid pgtype.UUID.
func PgUUIDFromGoogle(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: [16]byte(u), Valid: true}
}

// GoogleUUIDFromPg converts a valid pgtype.UUID to a google/uuid UUID.
//
// GoogleUUIDFromPg returns uuid.Nil and an error when u is not valid.
func GoogleUUIDFromPg(u pgtype.UUID) (uuid.UUID, error) {
	if !u.Valid {
		return uuid.Nil, errors.New("uuid is null")
	}
	return uuid.UUID(u.Bytes), nil
}

// PgUUIDFromString parses s as a UUID and returns it as a valid pgtype.UUID.
func PgUUIDFromString(s string) (pgtype.UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return PgUUIDFromGoogle(u), nil
}

// StringFromPgUUID converts a pgtype.UUID to its canonical string form.
//
// StringFromPgUUID returns an empty string when u is not valid.
func StringFromPgUUID(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return uuid.UUID(u.Bytes).String()
}

// TimeFromPgtype returns the time value stored in t.
func TimeFromPgtype(t pgtype.Timestamptz) time.Time {
	return t.Time
}

// TimePtrFromPgtype returns a pointer to t.Time when t is valid and non-zero.
//
// TimePtrFromPgtype returns nil for invalid or zero timestamps.
func TimePtrFromPgtype(t pgtype.Timestamptz) *time.Time {
	if !t.Valid || t.Time.IsZero() {
		return nil
	}
	return &t.Time
}

// Int64FromPgtype returns the int64 value stored in i.
func Int64FromPgtype(i pgtype.Int8) int64 {
	return i.Int64
}

// IntFromPgtype converts the int64 value stored in i to int.
func IntFromPgtype(i pgtype.Int8) int {
	return int(i.Int64)
}

// StringFromUUID converts a valid pgtype.UUID to a string.
//
// StringFromUUID returns an empty string when u is not valid.
func StringFromUUID(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return u.String()
}

// StringFromPgtype returns the string value stored in t.
//
// StringFromPgtype returns an empty string when t is not valid.
func StringFromPgtype(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// StringFromPoint formats a valid pgtype.Point as "x,y".
//
// StringFromPoint returns an empty string when p is not valid.
func StringFromPoint(p pgtype.Point) string {
	if !p.Valid {
		return ""
	}
	return strconv.FormatFloat(p.P.X, 'f', -1, 64) + "," + strconv.FormatFloat(p.P.Y, 'f', -1, 64)
}

// PgtypeFromTime converts t to a pgtype.Timestamptz.
//
// PgtypeFromTime returns an invalid timestamp when t is zero.
func PgtypeFromTime(t time.Time) pgtype.Timestamptz {
	if t.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// PgtypeFromInt64 converts i to a valid pgtype.Int8.
func PgtypeFromInt64(i int64) pgtype.Int8 {
	return pgtype.Int8{Int64: i, Valid: true}
}

// PgtypeFromInt converts i to a valid pgtype.Int8.
func PgtypeFromInt(i int) pgtype.Int8 {
	return pgtype.Int8{Int64: int64(i), Valid: true}
}

// PgtypeFromUUID scans s into a pgtype.UUID.
func PgtypeFromUUID(s string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(s)
	return uuid, err
}

// PgtypeFromString converts s to a valid pgtype.Text.
func PgtypeFromString(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

// PgtypeFromStringPtr converts s to pgtype.Text.
//
// PgtypeFromStringPtr returns an invalid text value when s is nil.
func PgtypeFromStringPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}
