// Package authorizationutil contains helpers for reading authorization data from
// request contexts.
package authorizationutil

import (
	"context"

	"github.com/google/uuid"
)

// AuthorizedResourceID returns the UUID stored in ctx under resource.
//
// The boolean result is false when the context value is missing, is not a
// uuid.UUID, or is uuid.Nil.
func AuthorizedResourceID(ctx context.Context, resource string) (uuid.UUID, bool) {
	resourceID, ok := ctx.Value(resource).(uuid.UUID)
	return resourceID, ok && resourceID != uuid.Nil
}
