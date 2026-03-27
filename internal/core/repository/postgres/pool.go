package postgres_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

	OpTimeout() time.Duration
}

type ConnectionPool struct {
	*pgxpool.Pool

	opTimeout time.Duration
}

func NewConnectionPool(ctx context.Context, cfg Config) (*ConnectionPool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Database,
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("pgxpool new: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("no database connection: %w", err)
	}

	return &ConnectionPool{
		Pool:      pool,
		opTimeout: cfg.Timeout,
	}, nil
}

func (p *ConnectionPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.Pool.QueryRow(ctx, sql, args...)
}

func (p *ConnectionPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.Pool.Query(ctx, sql, args...)
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.opTimeout
}
