package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type Repository[T any] interface {
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	GetAll(ctx context.Context) ([]*T, error)
	Create(ctx context.Context, domain *T) error
	Update(ctx context.Context, id uuid.UUID, task *T) error
}


type Services[T any] interface{
	
}