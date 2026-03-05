package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tclgroup/stock-management/internal/pkg/config"
)

// MustConnect creates a pgxpool and panics if the connection fails.
func MustConnect(cfg config.DatabaseConfig) *pgxpool.Pool {
	pool, err := Connect(cfg)
	if err != nil {
		panic(fmt.Sprintf("postgres: %v", err))
	}
	return pool
}

// Connect creates a pgxpool and verifies connectivity.
func Connect(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("new pool: %w", err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}

	return pool, nil
}
