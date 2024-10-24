package psql

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	psqldal "github.com/AAguilar0x0/txapp/extern/db/psql/dal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type DB struct {
	pool *pgxpool.Pool
	db   *psqldal.Queries
	tx   pgx.Tx
}

func New(host, user, password, db, port, sslmode string) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host,
		user,
		password,
		db,
		port,
		sslmode,
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

	err = dbInstance.Migrate("extern/db/psql/migrations", "up", nil, false)
	if err != nil {
		return nil, err
	}

	return &dbInstance, nil
}

func (d *DB) Migrate(dir, command string, version *int64, noVersioning bool) error {
	sqlDB := stdlib.OpenDBFromPool(d.pool)
	goose.SetBaseFS(os.DirFS("./"))
	if err := goose.SetDialect("pgx"); err != nil {
		return err
	}
	if strings.HasSuffix(command, "-to") && version == nil {
		return errors.New("missing version")
	}
	var err error
	options := []goose.OptionsFunc{}
	if noVersioning {
		options = append(options, goose.WithNoVersioning())
	}
	switch command {
	case "up":
		err = goose.Up(sqlDB, dir, options...)
	case "up-one":
		err = goose.UpByOne(sqlDB, dir, options...)
	case "up-to":
		err = goose.UpTo(sqlDB, dir, *version, options...)
	case "down":
		err = goose.Down(sqlDB, dir, options...)
	case "down-to":
		err = goose.DownTo(sqlDB, dir, *version, options...)
	default:
		return errors.New("missing or unknown migration command")
	}
	return err
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
