package psql

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func (d *Queries) Init(env services.Environment) error {
	if d.db != nil {
		return nil
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		env.GetDefault("DB_HOST", "localhost"),
		env.GetDefault("DB_USER", "postgres"),
		env.GetDefault("DB_PASSWORD", "postgres"),
		env.GetDefault("DB_NAME", "postgres"),
		env.GetDefault("DB_PORT", "5432"),
		env.GetDefault("DB_SSLMODE", "disable"),
	)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}
	config.MaxConns = int32(max(runtime.NumCPU(), 4))
	config.MinConns = 1
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}
	d.db = pool

	err = d.Migrate("extern/db/psql/migrations", "up", nil, false)
	if err != nil {
		return err
	}

	return nil
}

func (d *Queries) Close() {
	d.db.(*pgxpool.Pool).Close()
}

func (d *Queries) Migrate(dir, command string, version *int64, noVersioning bool) error {
	sqlDB := stdlib.OpenDBFromPool(d.db.(*pgxpool.Pool))
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

func (d *Queries) Begin(ctx context.Context) (*Queries, error) {
	val, ok := d.db.(*pgxpool.Pool)
	if d.db != nil && !ok {
		return d, nil
	}
	tx, err := val.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return New(tx), nil
}

func (d *Queries) Rollback(ctx context.Context) error {
	tx, ok := d.db.(pgx.Tx)
	if tx == nil || !ok {
		return nil
	}
	return tx.Rollback(ctx)
}

func (d *Queries) Commit(ctx context.Context) error {
	tx, ok := d.db.(pgx.Tx)
	if tx == nil || !ok {
		return nil
	}
	return tx.Commit(ctx)
}
