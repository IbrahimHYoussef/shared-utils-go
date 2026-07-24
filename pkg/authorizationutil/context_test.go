package authorizationutil

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestAuthorizedResourceID(t *testing.T) {
	resourceID := uuid.New()
	ctx := context.WithValue(context.Background(), "project", resourceID)

	got, ok := AuthorizedResourceID(ctx, "project")
	if !ok {
		t.Fatal("expected resource id to be authorized")
	}
	if got != resourceID {
		t.Fatalf("expected %s, got %s", resourceID, got)
	}
}

func TestAuthorizedResourceIDRejectsMissingAndNil(t *testing.T) {
	if _, ok := AuthorizedResourceID(context.Background(), "project"); ok {
		t.Fatal("expected missing resource id to be rejected")
	}

	ctx := context.WithValue(context.Background(), "project", uuid.Nil)
	if _, ok := AuthorizedResourceID(ctx, "project"); ok {
		t.Fatal("expected nil resource id to be rejected")
	}
}
