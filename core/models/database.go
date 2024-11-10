package models

import "context"

type Database interface {
	Migrate(dir, command string, version *int64, noVersioning bool) error
	Begin(ctx context.Context) (Database, error)
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
	Querier
}
