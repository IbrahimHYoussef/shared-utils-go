package dbutile

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

// fakeTx counts Commit and Rollback calls; every other pgx.Tx method is unused.
type fakeTx struct {
	pgx.Tx
	commits   int
	rollbacks int
}

func (f *fakeTx) Commit(context.Context) error   { f.commits++; return nil }
func (f *fakeTx) Rollback(context.Context) error { f.rollbacks++; return nil }

type fakeQueries struct{ tx pgx.Tx }

func (q fakeQueries) WithTx(tx pgx.Tx) fakeQueries { return fakeQueries{tx: tx} }

// ownerIn stores an owning transaction in ctx, as StartOrGet does when it
// begins one.
func ownerIn(ctx context.Context, tx *fakeTx) (context.Context, *Transaction[fakeQueries]) {
	owner := &Transaction[fakeQueries]{Qtx: fakeQueries{tx: tx}, Tx: tx, state: TxActive, txID: "tx_test"}
	return AddToCtx(ctx, owner), owner
}

func TestJoinedTransactionCannotFinalize(t *testing.T) {
	tx := &fakeTx{}
	ctx, owner := ownerIn(context.Background(), tx)

	// A nested service runs the usual StartOrGet / defer RollBack / Commit.
	joined, err := StartOrGet(ctx, nil, fakeQueries{})
	if err != nil {
		t.Fatalf("StartOrGet: %v", err)
	}
	if joined.Tx != tx || joined.Qtx.tx != tx {
		t.Fatal("joined handle must share the owner's pgx transaction and queries")
	}
	if err := joined.Commit(ctx); err != nil {
		t.Fatalf("joined Commit: %v", err)
	}
	if err := joined.RollBack(ctx); err != nil {
		t.Fatalf("joined RollBack: %v", err)
	}
	if tx.commits != 0 || tx.rollbacks != 0 {
		t.Fatalf("joined handle finalized the transaction: commits=%d rollbacks=%d", tx.commits, tx.rollbacks)
	}

	if err := owner.Commit(ctx); err != nil {
		t.Fatalf("owner Commit: %v", err)
	}
	_ = owner.RollBack(ctx) // deferred rollback after commit stays a no-op
	if tx.commits != 1 || tx.rollbacks != 0 {
		t.Fatalf("owner must commit exactly once: commits=%d rollbacks=%d", tx.commits, tx.rollbacks)
	}
}

func TestJoinedFailureLeavesRollbackToOwner(t *testing.T) {
	tx := &fakeTx{}
	ctx, owner := ownerIn(context.Background(), tx)

	joined, _ := StartOrGet(ctx, nil, fakeQueries{})
	_ = joined.RollBack(ctx) // nested workflow failed and ran its deferred rollback
	if tx.rollbacks != 0 {
		t.Fatal("joined handle rolled back the owner's transaction")
	}

	_ = owner.RollBack(ctx) // the error propagates; the owner rolls back
	if tx.rollbacks != 1 || tx.commits != 0 {
		t.Fatalf("owner must roll back exactly once: commits=%d rollbacks=%d", tx.commits, tx.rollbacks)
	}
}

func TestAddToCtxKeepsTheOwner(t *testing.T) {
	tx := &fakeTx{}
	ctx, owner := ownerIn(context.Background(), tx)
	joined, _ := StartOrGet(ctx, nil, fakeQueries{})

	if got := AddToCtx(ctx, joined); got != ctx {
		t.Fatal("AddToCtx must not replace the owner with a joined handle")
	}
	if stored, _ := GetFromCtx[fakeQueries](ctx); stored != owner {
		t.Fatal("context must still hold the owning transaction")
	}
}
