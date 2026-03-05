package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	pg "github.com/tclgroup/stock-management/internal/adapter/postgres"
)

// querier abstracts over *pgxpool.Pool and pgx.Tx so repos can use either.
type querier interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

// getQuerier returns the transaction from ctx if present, otherwise the pool.
func getQuerier(ctx context.Context, pool *pgxpool.Pool) querier {
	if tx, ok := pg.ExtractTx(ctx); ok {
		return tx
	}
	return pool
}
