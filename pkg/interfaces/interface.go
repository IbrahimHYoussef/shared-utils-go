// Package interfaces contains generic application-layer interfaces shared across
// services and repositories.
package interfaces

import (
	"context"

	"github.com/google/uuid"
)

// Repository describes the common persistence operations for a domain type T.
type Repository[T any] interface {
	// GetByID returns a domain object by id.
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	// GetAll returns all domain objects.
	GetAll(ctx context.Context) ([]*T, error)
	// Create stores a new domain object.
	Create(ctx context.Context, domain *T) error
	// Update replaces or updates the domain object identified by id.
	Update(ctx context.Context, id uuid.UUID, task *T) error
}

// Services is a marker interface for shared service abstractions over T.
type Services[T any] interface {
}
