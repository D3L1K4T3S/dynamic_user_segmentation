package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"

	e "dynamic-user-segmentation/pkg/util/erros"
)

const (
	defaultConnectionTimeout  = 1 * time.Second
	defaultMaxPoolSize        = 1
	defaultConnectionAttempts = 10
)

type PgxPool interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}

type PostgreSQL struct {
	Pool PgxPool

	connTimeout  time.Duration
	connAttempts int
	maxPoolSize  int
}

func NewClient(url string, options ...Option) (*PostgreSQL, error) {

	pg := &PostgreSQL{
		connTimeout:  defaultConnectionTimeout,
		connAttempts: defaultConnectionAttempts,
		maxPoolSize:  defaultMaxPoolSize,
	}

	for _, option := range options {
		option(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, e.Wrap("can't parse config: ", err)
	}

	defer func() {
		err = e.WrapIfErr("Can't crate a postgres client", err)
	}()

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Trying to connect, attempts left: %d", pg.connAttempts)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, err
	}

	return pg, nil
}

func (p *PostgreSQL) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
