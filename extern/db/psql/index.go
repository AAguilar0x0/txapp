package psql

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/AAguilar0x0/txapp/extern/db/psql/dal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Psql struct {
	pool *pgxpool.Pool
	db   *dal.Queries
	tx   pgx.Tx
}

func transformError(err error) *apierrors.APIError {
	if errors.Is(err, pgx.ErrNoRows) {
		return apierrors.NotFound(err.Error())
	}
	return apierrors.InternalServerError(err.Error())
}

func New(env services.Environment) (*Psql, error) {
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
	a := Psql{
		pool: pool,
	}
	d := dal.New(pool)
	a.db = d

	err = a.Migrate("extern/db/psql/migrations", "up", nil, false)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (d *Psql) Close() error {
	d.pool.Close()
	return nil
}

func (d *Psql) Migrate(dir, command string, version *int64, noVersioning bool) error {
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

func (d *Psql) Begin(ctx context.Context) (models.Database, error) {
	if d.tx != nil {
		return d, nil
	}
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &Psql{
		tx: tx,
		db: dal.New(tx),
	}, nil
}

func (d *Psql) Rollback(ctx context.Context) error {
	if d.tx == nil {
		return nil
	}
	return d.tx.Rollback(ctx)
}

func (d *Psql) Commit(ctx context.Context) error {
	if d.tx == nil {
		return nil
	}
	return d.tx.Commit(ctx)
}
