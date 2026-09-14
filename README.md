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

### dbutile
PostgreSQL connection, migration, transaction, and repository helpers built on `pgx`.

`BaseRepository` is intended for repositories backed by sqlc-generated query types. The query type must implement `WithTx(pgx.Tx) Q`, which sqlc commonly generates for pgx projects.

```go
func (q *Queries) WithTx(tx pgx.Tx) *Queries
```

```go
type AuthRepository struct {
    dbutile.BaseRepository[*database.Queries]
}

func NewAuthRepository(queries *database.Queries) AuthRepository {
    return AuthRepository{
        BaseRepository: dbutile.NewBaseRepository(queries),
    }
}

func (r AuthRepository) FindUserByEmail(ctx context.Context, email string) (auth.User, error) {
    row, err := r.QTX(ctx).GetUserByEmail(ctx, email)
    if err != nil {
        return auth.User{}, err
    }

    return authUserFromDatabase(row), nil
}
```

If you prefer a local lowercase helper, wrap `QTX` in your concrete repository:

```go
func (r AuthRepository) qtx(ctx context.Context) *database.Queries {
    return r.QTX(ctx)
}
```

Use the same query object when starting transactions. Any repository method that calls `QTX(ctx)` will automatically use the transaction-bound queries after the transaction is added to the context.

```go
func RunInTransaction(ctx context.Context, pool *pgxpool.Pool, queries *database.Queries, repo AuthRepository) error {
    tx, err := dbutile.StartOrGet(ctx, pool, queries)
    if err != nil {
        return err
    }

    ctx = dbutile.AddToCtx(ctx, tx)
    defer tx.Rollback(ctx)

    if _, err := repo.FindUserByEmail(ctx, "user@example.com"); err != nil {
        return err
    }

    return tx.Commit(ctx)
}
```

## Installation

```bash
go get github.com/yourusername/project-management/shared-utils-go
```
