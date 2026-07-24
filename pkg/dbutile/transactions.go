package dbutile

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionKeyType struct{}

var transactionKey = transactionKeyType{}

type TxState int

const (
	TxActive TxState = iota
	TxCommitted
	TxRolledBack
)

type TransactionalQueries[Q any] interface {
	WithTx(pgx.Tx) Q
}

type Transaction[Q TransactionalQueries[Q]] struct {
	Qtx       Q
	Tx        pgx.Tx
	state     TxState
	startTime time.Time
	txID      string
}

func StartOrGet[Q TransactionalQueries[Q]](ctx context.Context, db *pgxpool.Pool, qtx Q) (*Transaction[Q], error) {
	t, exists := GetFromCtx[Q](ctx)
	if exists {
		logger.Debug("Using existing transaction from context", "tx_id", t.txID)
		return t, nil
	}

	txID := fmt.Sprintf("tx_%d", time.Now().UnixNano())
	startTime := time.Now()

	logger.Info("Starting new database transaction", "tx_id", txID)

	tx, err := db.Begin(ctx)
	if err != nil {
		logger.Error("Error beginning transaction",
			"error", err,
			"tx_id", txID,
			"duration_ms", time.Since(startTime).Milliseconds(),
		)
		return nil, errors.New("Internal server error")
	}

	transaction := &Transaction[Q]{
		Qtx:       qtx.WithTx(tx),
		Tx:        tx,
		state:     TxActive,
		startTime: startTime,
		txID:      txID,
	}

	logger.Debug("Database transaction started successfully",
		"tx_id", txID,
		"duration_ms", time.Since(startTime).Milliseconds(),
	)

	return transaction, nil
}

func (t *Transaction[Q]) RollBack(ctx context.Context) error {
	if t.state != TxActive {
		logger.Debug("Transaction already finalized, skipping rollback",
			"tx_id", t.txID,
			"current_state", t.state,
		)
		return nil
	}

	duration := time.Since(t.startTime)

	logger.Info("Rolling back database transaction",
		"tx_id", t.txID,
		"duration_ms", duration.Milliseconds(),
	)

	err := t.Tx.Rollback(ctx)
	if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		logger.Error("Error rolling back transaction",
			"error", err,
			"tx_id", t.txID,
			"duration_ms", duration.Milliseconds(),
		)
		return err
	}

	t.state = TxRolledBack

	logger.Info("Database transaction rolled back successfully",
		"tx_id", t.txID,
		"duration_ms", duration.Milliseconds(),
	)

	return nil
}

func (t *Transaction[Q]) Rollback(ctx context.Context) error {
	return t.RollBack(ctx)
}

func (t *Transaction[Q]) Commit(ctx context.Context) error {
	duration := time.Since(t.startTime)

	logger.Info("Committing database transaction",
		"tx_id", t.txID,
		"duration_ms", duration.Milliseconds(),
	)

	err := t.Tx.Commit(ctx)
	if err != nil {
		logger.Error("Error committing transaction",
			"error", err,
			"tx_id", t.txID,
			"duration_ms", duration.Milliseconds(),
		)
		return err
	}

	t.state = TxCommitted

	logger.Info("Database transaction committed successfully",
		"tx_id", t.txID,
		"duration_ms", duration.Milliseconds(),
	)

	return nil
}

func AddToCtx[Q TransactionalQueries[Q]](ctx context.Context, t *Transaction[Q]) context.Context {
	if _, exists := GetFromCtx[Q](ctx); exists {
		return ctx
	}

	return context.WithValue(ctx, transactionKey, t)
}

func GetFromCtx[Q TransactionalQueries[Q]](ctx context.Context) (*Transaction[Q], bool) {
	tx, ok := ctx.Value(transactionKey).(*Transaction[Q])
	return tx, ok
}
