package authorizationutil

import (
	"context"

	"github.com/google/uuid"
)

func AuthorizedResourceID(ctx context.Context, resource string) (uuid.UUID, bool) {
	resourceID, ok := ctx.Value(resource).(uuid.UUID)
	return resourceID, ok && resourceID != uuid.Nil
}
