package mapping

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// StringPtrFromPgtype returns a pointer to t.String when t is valid.
//
// StringPtrFromPgtype returns nil when t is not valid.
func StringPtrFromPgtype(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

// Int64PtrFromPgtype returns a pointer to i.Int64 when i is valid.
//
// Int64PtrFromPgtype returns nil when i is not valid.
func Int64PtrFromPgtype(i pgtype.Int8) *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// PgtypeFromInt64Ptr converts i to pgtype.Int8.
//
// PgtypeFromInt64Ptr returns an invalid int8 value when i is nil.
func PgtypeFromInt64Ptr(i *int64) pgtype.Int8 {
	if i == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *i, Valid: true}
}

// PgtypeFromTimePtr converts t to pgtype.Timestamptz.
//
// PgtypeFromTimePtr returns an invalid timestamp when t is nil or zero.
func PgtypeFromTimePtr(t *time.Time) pgtype.Timestamptz {
	if t == nil || t.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// BytesFromStringPtr converts s to a byte slice, for json and jsonb columns
// that sqlc maps to []byte.
//
// BytesFromStringPtr returns nil when s is nil, which pgx writes as NULL.
func BytesFromStringPtr(s *string) []byte {
	if s == nil {
		return nil
	}
	return []byte(*s)
}

// StringPtrFromBytes converts b to a string pointer, for json and jsonb
// columns that sqlc maps to []byte.
//
// StringPtrFromBytes returns nil when b is empty.
func StringPtrFromBytes(b []byte) *string {
	if len(b) == 0 {
		return nil
	}
	s := string(b)
	return &s
}
