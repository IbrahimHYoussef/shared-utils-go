# Shared Go Utilities

Common Go utility functions used across Go services.

## Packages

### jwt
JWT token generation and validation utilities.

```go
import "github.com/yourusername/project-management/shared-utils-go/jwt"

tm := jwt.NewTokenManager("your-secret-key")
token, err := tm.GenerateToken("user123", "user@example.com")
claims, err := tm.ValidateToken(token)
```

### crypto
Password hashing and verification using bcrypt.

```go
import "github.com/yourusername/project-management/shared-utils-go/crypto"

hash, err := crypto.HashPassword("mypassword")
isValid := crypto.CheckPasswordHash("mypassword", hash)
```

### validator
Input validation utilities.

```go
import "github.com/yourusername/project-management/shared-utils-go/validator"

if validator.IsValidEmail("user@example.com") {
    // email is valid
}

if validator.IsValidPassword("securepass123") {
    // password meets requirements
}
```

## Installation

```bash
go get github.com/yourusername/project-management/shared-utils-go
```
