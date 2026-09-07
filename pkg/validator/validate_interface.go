package validator

import "github.com/google/uuid"

// Validator describes a type that can validate itself.
type Validator interface {
	// IsValid returns nil when the value is valid.
	IsValid() error
}

// Storable describes a type that can convert itself to a database model.
type Storable[D any] interface {
	// ToDB converts the value into its database representation.
	ToDB() (D, error)
}

// StorableByID describes a type that can convert itself to a database model
// using an external UUID.
type StorableByID[D Validator] interface {
	// ToDB converts the value into its database representation using id.
	ToDB(uuid.UUID) (D, error)
}

// Sendable describes a type that can convert itself into a DTO.
type Sendable[T Validator] interface {
	// ToDto converts the value into its DTO representation.
	ToDto() (T, error)
}
