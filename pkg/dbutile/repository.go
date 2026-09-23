// Package dbutile contains PostgreSQL connection, migration, transaction, and
// logging helpers built on pgx.
package dbutile

import (
	"context"
)

// BaseRepository stores a default query object and resolves transaction-bound
// queries from context for repository methods.
type BaseRepository[Q TransactionalQueries[Q]] struct {
	queries Q
}

// NewBaseRepository creates a BaseRepository with the default non-transactional
// query object.
func NewBaseRepository[Q TransactionalQueries[Q]](queries Q) BaseRepository[Q] {
	return BaseRepository[Q]{queries: queries}
}

// QTX returns the transaction-bound query object from ctx when a transaction is
// active, otherwise it returns the repository's default query object. After
// the transaction in ctx commits or rolls back, QTX uses the default query
// object again.
func (r BaseRepository[Q]) QTX(ctx context.Context) Q {
	if trans, ok := GetFromCtx[Q](ctx); ok {
		return trans.Qtx
	}

	return r.queries
}
