package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

// TxManager implements port.TxManager using pgxpool transactions.
type TxManager struct {
	pool *pgxpool.Pool
}

// NewTxManager creates a new TxManager.
func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

// WithinTx executes fn inside a database transaction.
// The transaction is stored on the context so repositories can use it.
func (m *TxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err = fn(txCtx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("rollback: %v, original: %w", rbErr, err)
		}
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// ExtractTx retrieves a pgx.Tx from the context, if present.
func ExtractTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}
