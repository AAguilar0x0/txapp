package psql

import (
	"context"
	"fmt"
	"runtime"
	"time"

	psqldal "github.com/AAguilar0x0/bapp/extern/db/psql/dal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
	db   *psqldal.Queries
	tx   pgx.Tx
}

func New(host, user, password, db, port string) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=verify-full",
		host,
		user,
		password,
		db,
		port,
	)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	config.MaxConns = int32(max(runtime.NumCPU(), 4))
	config.MinConns = 1
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	dbInstance := DB{
		pool: pool,
		db:   psqldal.New(pool),
	}
	return &dbInstance, nil
}

func (d *DB) Close() {
	d.pool.Close()
}

func (d *DB) Instance() *psqldal.Queries {
	return d.db
}

func (d *DB) Begin(ctx context.Context) (*DB, error) {
	if d.tx != nil {
		return d, nil
	}
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &DB{
		tx: tx,
		db: d.db.WithTx(tx),
	}, nil
}

func (d *DB) Rollback(ctx context.Context) error {
	if d.tx == nil {
		return nil
	}
	return d.tx.Rollback(ctx)
}

func (d *DB) Commit(ctx context.Context) error {
	if d.tx == nil {
		return nil
	}
	return d.tx.Commit(ctx)
}