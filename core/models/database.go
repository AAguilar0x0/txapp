package models

import (
	"context"
	"io/fs"

	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
)

type Database interface {
	Begin(ctx context.Context) (Database, *apierrors.APIError)
	Rollback(ctx context.Context) *apierrors.APIError
	Commit(ctx context.Context) *apierrors.APIError
	Querier
}

type Migrator interface {
	Migrate(fsys fs.FS, dir, command string, version *int64, noVersioning bool) *apierrors.APIError
}

type DatabaseManager interface {
	Database
	Migrator
}
