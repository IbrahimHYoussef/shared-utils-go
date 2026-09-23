package mapping

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestStringPtrFromPgtype(t *testing.T) {
	if got := StringPtrFromPgtype(pgtype.Text{}); got != nil {
		t.Fatalf("invalid text: got %v, want nil", *got)
	}
	got := StringPtrFromPgtype(pgtype.Text{String: "a", Valid: true})
	if got == nil || *got != "a" {
		t.Fatalf("valid text: got %v, want a", got)
	}
	empty := StringPtrFromPgtype(pgtype.Text{String: "", Valid: true})
	if empty == nil || *empty != "" {
		t.Fatalf("valid empty text: got %v, want pointer to empty string", empty)
	}
}

func TestInt64PtrRoundTrip(t *testing.T) {
	if got := Int64PtrFromPgtype(pgtype.Int8{}); got != nil {
		t.Fatalf("invalid int8: got %v, want nil", *got)
	}
	if got := PgtypeFromInt64Ptr(nil); got.Valid {
		t.Fatalf("nil pointer: got valid int8")
	}
	v := int64(42)
	got := Int64PtrFromPgtype(PgtypeFromInt64Ptr(&v))
	if got == nil || *got != 42 {
		t.Fatalf("round trip: got %v, want 42", got)
	}
}

func TestPgtypeFromTimePtr(t *testing.T) {
	if got := PgtypeFromTimePtr(nil); got.Valid {
		t.Fatalf("nil pointer: got valid timestamp")
	}
	zero := time.Time{}
	if got := PgtypeFromTimePtr(&zero); got.Valid {
		t.Fatalf("zero time: got valid timestamp")
	}
	now := time.Now()
	got := PgtypeFromTimePtr(&now)
	if !got.Valid || !got.Time.Equal(now) {
		t.Fatalf("valid time: got %v, want %v", got, now)
	}
}

func TestBytesStringPtrRoundTrip(t *testing.T) {
	if got := BytesFromStringPtr(nil); got != nil {
		t.Fatalf("nil pointer: got %v, want nil", got)
	}
	if got := StringPtrFromBytes(nil); got != nil {
		t.Fatalf("nil bytes: got %v, want nil", *got)
	}
	s := `{"k":"v"}`
	got := StringPtrFromBytes(BytesFromStringPtr(&s))
	if got == nil || *got != s {
		t.Fatalf("round trip: got %v, want %s", got, s)
	}
}
