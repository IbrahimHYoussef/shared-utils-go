package mapping

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UUID Covnerstions
func PgUUIDFromGoogle(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: [16]byte(u), Valid: true}
}
func GoogleUUIDFromPg(u pgtype.UUID) (uuid.UUID, error) {
	if !u.Valid {
		return uuid.Nil, errors.New("uuid is null")
	}
	return uuid.UUID(u.Bytes), nil
}
func PgUUIDFromString(s string) (pgtype.UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return PgUUIDFromGoogle(u), nil
}
func StringFromPgUUID(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return uuid.UUID(u.Bytes).String()
}

// Convert pgtype to standard Go types
func TimeFromPgtype(t pgtype.Timestamptz) time.Time {
	return t.Time
}

func TimePtrFromPgtype(t pgtype.Timestamptz) *time.Time {
	if !t.Valid || t.Time.IsZero() {
		return nil
	}
	return &t.Time
}

func Int64FromPgtype(i pgtype.Int8) int64 {
	return i.Int64
}

func IntFromPgtype(i pgtype.Int8) int {
	return int(i.Int64)
}

func StringFromUUID(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return u.String()
}

func StringFromPgtype(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func StringFromPoint(p pgtype.Point) string {
	if !p.Valid {
		return ""
	}
	return strconv.FormatFloat(p.P.X, 'f', -1, 64) + "," + strconv.FormatFloat(p.P.Y, 'f', -1, 64)
}

// Reverse conversions
func PgtypeFromTime(t time.Time) pgtype.Timestamptz {
	if t.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func PgtypeFromInt64(i int64) pgtype.Int8 {
	return pgtype.Int8{Int64: i, Valid: true}
}

func PgtypeFromInt(i int) pgtype.Int8 {
	return pgtype.Int8{Int64: int64(i), Valid: true}
}

func PgtypeFromUUID(s string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(s)
	return uuid, err
}

func PgtypeFromString(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

func PgtypeFromStringPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}
