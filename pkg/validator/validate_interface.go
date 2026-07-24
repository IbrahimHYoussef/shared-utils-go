package validator

import "github.com/google/uuid"

type Validator interface {
	IsValid() error
}

type Storable[D any] interface {
	ToDB() (D, error)
}

type StorableByID[D Validator] interface {
	ToDB(uuid.UUID) (D, error)
}

type Sendable[T Validator] interface {
	ToDto() (T, error)
}
