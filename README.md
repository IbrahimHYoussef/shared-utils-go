# Shared Go Utilities

Common Go utility functions used across Go services.

## Packages

### jwtutils
HS256 JWT access tokens and random refresh session tokens.

```go
import "github.com/IbrahimHYoussef/shared-utils-go/pkg/jwtutils"

jwtService := jwtutils.NewJwtManager(secret, 900, 2592000, "my-service") // access s, refresh s, issuer
token, err := jwtService.GenerateToken(userID, "user@example.com", "user")
claims, err := jwtService.ValidateToken(token)
refresh, err := jwtService.GenerateRefreshSessionToken(32)
```

`ValidateToken` and `ParseWithClaims` accept only tokens that are signed with HS256, carry an `exp` claim, and have an `iss` claim equal to the service issuer (the issuer check is skipped when the issuer is empty). `ParseOptions(issuer)` returns the same rules for custom parsing.

### middelware (authentication)
`AuthMiddleWareFactoryFromService(jwtService)` authenticates `Authorization: Bearer <token>` requests with the `jwtutils` rules above and stores `*jwtutils.Claims` in the context under `UserClaimsKey`. Prefer it over `AuthMiddleWareFactory(secret)`, which applies the same rules but cannot check the issuer. Expired tokens get a `Token Expired` 401; every other failure gets `Not Allowed To Access This Endpoint`.

### crypto
Password hashing and verification using bcrypt.

```go
import "github.com/yourusername/project-management/shared-utils-go/crypto"

hash, err := crypto.HashPassword("mypassword")
isValid := crypto.CheckPasswordHash("mypassword", hash)
```

For high-entropy secrets that are stored hashed (refresh tokens, emailed link tokens, OTP codes), use the SHA-256 helpers instead of bcrypt:

```go
token, err := crypto.GenerateSecretToken(32) // 32 random bytes, unpadded base64url
stored := crypto.HashToken(token)            // lowercase hex SHA-256
ok := crypto.CheckTokenHash(token, stored)   // constant-time comparison
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

Only the call that begins a transaction owns it. When `StartOrGet` finds a transaction already in the context, it returns a joined handle to the same transaction whose `Commit` and `RollBack` do nothing, so a service called inside another service's workflow can keep its usual `StartOrGet` / `defer RollBack` / `Commit` code without finalizing the caller's work. The owner commits once everything succeeded, or rolls back when an error propagates up.

A `Transaction` is a small state machine: `Active` until its owner commits or rolls back (a failed commit also ends it), then `Committed` or `RolledBack`; `State()` reports it. A finished transaction stays in its context, but `GetFromCtx` only returns active ones, so code that keeps using that context after the commit — `QTX` in repositories, or a later `StartOrGet` — gets the default queries or a new transaction instead of a closed one. Work done after a commit (sending an email, for example) is therefore safe with the same `ctx`.

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

Use `MapError` and the unique-violation helpers to translate pgx errors inside repositories, so services never import pgx:

```go
row, err := r.QTX(ctx).CreateUser(ctx, params)
if dbutile.IsUniqueViolation(err) {
    return auth.User{}, auth.ErrEmailOrUsernameTaken
}
if err != nil {
    return auth.User{}, dbutile.MapError(err) // pgx.ErrNoRows -> dbutile.ErrNotFound
}
```

`UniqueViolationConstraint(err)` also returns the violated constraint name, for example `users_email_lower_key`.

### mapping
Conversion helpers between Go types and pgx `pgtype` values, used inside repository implementations.

| Go type | pgtype | To pgtype | From pgtype |
| --- | --- | --- | --- |
| `string` | `pgtype.Text` | `PgtypeFromString` | `StringFromPgtype` |
| `*string` | `pgtype.Text` | `PgtypeFromStringPtr` | `StringPtrFromPgtype` |
| `int64` / `int` | `pgtype.Int8` | `PgtypeFromInt64`, `PgtypeFromInt` | `Int64FromPgtype`, `IntFromPgtype` |
| `*int64` | `pgtype.Int8` | `PgtypeFromInt64Ptr` | `Int64PtrFromPgtype` |
| `time.Time` | `pgtype.Timestamptz` | `PgtypeFromTime` | `TimeFromPgtype` |
| `*time.Time` | `pgtype.Timestamptz` | `PgtypeFromTimePtr` | `TimePtrFromPgtype` |
| `uuid.UUID` | `pgtype.UUID` | `PgUUIDFromGoogle` | `GoogleUUIDFromPg` |
| `*string` | `[]byte` (json/jsonb) | `BytesFromStringPtr` | `StringPtrFromBytes` |

### testutils
Postgres testcontainers for integration tests (Docker required).

```go
container, err := testutils.CreatePGContainer(ctx, "", nil) // postgres:18-alpine, db/user/password "test"
defer container.Terminate(ctx)
pool, err := pgxpool.New(ctx, container.ConnectionString)
```

Pass a `*PostgresTestCred` to choose the database, user or password; empty fields fall back to `test`. The container is ready once Postgres accepts connections (up to 60 s).

## Installation

```bash
go get github.com/yourusername/project-management/shared-utils-go
```
