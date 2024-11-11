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
	pool  *pgxpool.Pool
	db    *dal.Queries
	tx    pgx.Tx
	idGen services.IDGenerator
}

func transformError(err error) *apierrors.APIError {
	if errors.Is(err, pgx.ErrNoRows) {
		return apierrors.NotFound("not found", err.Error())
	}
	return apierrors.InternalServerError("an error occurred", err.Error())
}

func New(env services.Environment, idGen services.IDGenerator) (*Psql, error) {
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
		pool:  pool,
		idGen: idGen,
	}
	d := dal.New(pool)
	a.db = d

	apiErr := a.Migrate("extern/db/psql/migrations", "up", nil, false)
	if apiErr != nil {
		return nil, apiErr
	}

	return &a, nil
}

func (d *Psql) Close() error {
	d.pool.Close()
	return nil
}

func (d *Psql) Migrate(dir, command string, version *int64, noVersioning bool) *apierrors.APIError {
	sqlDB := stdlib.OpenDBFromPool(d.pool)
	goose.SetBaseFS(os.DirFS("./"))
	if err := goose.SetDialect("pgx"); err != nil {
		return apierrors.InternalServerError(err.Error())
	}
	if strings.HasSuffix(command, "-to") && version == nil {
		return apierrors.InternalServerError("missing version")
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
		return apierrors.InternalServerError("missing or unknown migration command")
	}
	if err != nil {
		return apierrors.InternalServerError(err.Error())
	}
	return nil
}

func (d *Psql) Begin(ctx context.Context) (models.Database, *apierrors.APIError) {
	if d.tx != nil {
		return d, nil
	}
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return nil, apierrors.InternalServerError("cannot create transaction", err.Error())
	}
	return &Psql{
		tx:    tx,
		db:    dal.New(tx),
		idGen: d.idGen,
	}, nil
}

func (d *Psql) Rollback(ctx context.Context) *apierrors.APIError {
	if d.tx == nil {
		return nil
	}
	err := d.tx.Rollback(ctx)
	if err != nil {
		return apierrors.InternalServerError("cannot rollback", err.Error())
	}
	return nil
}

func (d *Psql) Commit(ctx context.Context) *apierrors.APIError {
	if d.tx == nil {
		return nil
	}
	err := d.tx.Commit(ctx)
	if err != nil {
		return apierrors.InternalServerError("cannot commit", err.Error())
	}
	return nil
}
