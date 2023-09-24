// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/driftprogramming/pgxpoolmock"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    pgxpoolmock.PgxPool
}

// New -.
func New(ctx context.Context, url string, opts ...Option) (*Postgres, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("postgres - New - ctx.Done: %w", ctx.Err())
	default:
		postgres := &Postgres{
			maxPoolSize:  _defaultMaxPoolSize,
			connAttempts: _defaultConnAttempts,
			connTimeout:  _defaultConnTimeout,
		}

		// Custom options
		for _, opt := range opts {
			opt(postgres)
		}

		postgres.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		poolConfig, err := pgxpool.ParseConfig(url)
		if err != nil {
			return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
		}

		poolConfig.MaxConns = int32(postgres.maxPoolSize)

		for postgres.connAttempts > 0 {
			postgres.Pool, err = pgxpool.ConnectConfig(ctx, poolConfig)
			if err == nil {
				break
			}

			log.Printf("Postgres is trying to connect, attempts left: %d", postgres.connAttempts)

			time.Sleep(postgres.connTimeout)

			postgres.connAttempts--
		}

		if err != nil {
			return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
		}

		return postgres, nil
	}
}

// Close -.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
