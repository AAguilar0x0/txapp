package models

import (
	"context"

	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
)

type Database interface {
	Migrate(dir, command string, version *int64, noVersioning bool) *apierrors.APIError
	Begin(ctx context.Context) (Database, *apierrors.APIError)
	Rollback(ctx context.Context) *apierrors.APIError
	Commit(ctx context.Context) *apierrors.APIError
	Querier
}
