package paging_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/IbrahimHYoussef/shared-utils-go/pkg/paging"
	"github.com/google/uuid"
)

func TestClampLimit(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want int
	}{
		{"absent", "", 30},
		{"zero", "0", 30},
		{"not a number", "abc", 30},
		{"NaN", "NaN", 30},
		{"fraction below one", "0.9", 30},
		{"in range", "10", 10},
		{"spaces", " 10 ", 10},
		{"fraction rounded down", "10.9", 10},
		{"one", "1", 1},
		{"negative", "-5", 1},
		{"negative fraction", "-0.5", 1},
		{"at maximum", "100", 100},
		{"above maximum", "101", 100},
		{"far above maximum", "1e30", 100},
		{"infinity", "Inf", 100},
		{"negative infinity", "-Inf", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := paging.ClampLimit(tt.raw, 30, 100); got != tt.want {
				t.Fatalf("ClampLimit(%q, 30, 100): got %d, want %d", tt.raw, got, tt.want)
			}
		})
	}
}

func TestCursorRoundTrip(t *testing.T) {
	want := paging.Cursor[time.Time]{
		Key: time.Date(2026, 10, 1, 9, 46, 0, 0, time.UTC),
		ID:  uuid.MustParse("0199a3f1-6c80-7d21-a9e3-1f4b5c6d7e03"),
	}

	cursor, err := paging.EncodeCursor(want)
	if err != nil {
		t.Fatalf("EncodeCursor: %s", err)
	}
	got, err := paging.DecodeCursor[paging.Cursor[time.Time]](cursor)
	if err != nil {
		t.Fatalf("DecodeCursor: %s", err)
	}
	if got == nil || !got.Key.Equal(want.Key) || got.ID != want.ID {
		t.Fatalf("round trip: got %+v, want %+v", got, want)
	}
}

func TestDecodeCursorEmpty(t *testing.T) {
	got, err := paging.DecodeCursor[paging.Cursor[string]]("")
	if err != nil {
		t.Fatalf("empty cursor: got error %s, want nil", err)
	}
	if got != nil {
		t.Fatalf("empty cursor: got %+v, want nil", got)
	}
}

func TestDecodeCursorInvalid(t *testing.T) {
	wrongType, err := paging.EncodeCursor(paging.Cursor[string]{Key: "blue", ID: uuid.New()})
	if err != nil {
		t.Fatalf("EncodeCursor: %s", err)
	}

	tests := []struct {
		name   string
		cursor string
	}{
		{"not base64url", "!!!"},
		{"not json", "bm90IGpzb24"},
		{"numeric offset", "30"},
		{"another key type", wrongType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := paging.DecodeCursor[paging.Cursor[time.Time]](tt.cursor)
			if !errors.Is(err, paging.ErrInvalidCursor) {
				t.Fatalf("got error %v, want ErrInvalidCursor", err)
			}
			if got != nil {
				t.Fatalf("got %+v, want nil", got)
			}
		})
	}
}

type row struct {
	ID   uuid.UUID
	Name string
}

func rowKey(r row) paging.Cursor[string] {
	return paging.Cursor[string]{Key: r.Name, ID: r.ID}
}

func TestNewPageHasNext(t *testing.T) {
	rows := []row{
		{uuid.New(), "a"},
		{uuid.New(), "b"},
		{uuid.New(), "c"},
	}

	page, err := paging.NewPage(rows, 2, rowKey)
	if err != nil {
		t.Fatalf("NewPage: %s", err)
	}
	if len(page.Items) != 2 || page.Items[1].Name != "b" {
		t.Fatalf("items: got %+v, want the first two rows", page.Items)
	}
	if page.NextCursor == nil {
		t.Fatalf("next cursor: got nil, want a cursor")
	}
	key, err := paging.DecodeCursor[paging.Cursor[string]](*page.NextCursor)
	if err != nil {
		t.Fatalf("DecodeCursor: %s", err)
	}
	if key.Key != "b" || key.ID != rows[1].ID {
		t.Fatalf("next cursor: got %+v, want the last kept row", key)
	}
}

func TestNewPageLastPage(t *testing.T) {
	rows := []row{{uuid.New(), "a"}, {uuid.New(), "b"}}

	page, err := paging.NewPage(rows, 2, rowKey)
	if err != nil {
		t.Fatalf("NewPage: %s", err)
	}
	if len(page.Items) != 2 {
		t.Fatalf("items: got %d, want 2", len(page.Items))
	}
	if page.NextCursor != nil {
		t.Fatalf("next cursor: got %q, want nil", *page.NextCursor)
	}
}

func TestNewPageEmptyJSON(t *testing.T) {
	page, err := paging.NewPage[row](nil, 10, rowKey)
	if err != nil {
		t.Fatalf("NewPage: %s", err)
	}
	got, err := json.Marshal(page)
	if err != nil {
		t.Fatalf("Marshal: %s", err)
	}
	if want := `{"items":[],"next_cursor":null}`; string(got) != want {
		t.Fatalf("json: got %s, want %s", got, want)
	}
}
