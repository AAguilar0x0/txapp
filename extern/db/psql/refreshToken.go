package psql

import (
	"context"
	"time"

	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/pkg/withretry"
)

func (d *Psql) RefreshTokenCreate(ctx context.Context, id, userID string, expiration time.Time) *apierrors.APIError {
	errI := withretry.WithRetry(ctx, withretry.DefaultConfig, transientError, func(ctx context.Context) error {
		return d.db.RefreshTokenCreate(ctx, id, userID, expiration)
	})
	return transformError(errI)
}

func (d *Psql) RefreshTokenGet(ctx context.Context, id, userID string) (*models.RefreshToken, *apierrors.APIError) {
	var data *models.RefreshToken
	errI := withretry.WithRetry(ctx, withretry.DefaultConfig, transientError, func(ctx context.Context) error {
		result, err := d.db.RefreshTokenGet(ctx, id, userID)
		data = (*models.RefreshToken)(result)
		return err
	})
	return data, transformError(errI)
}

func (d *Psql) RefreshTokenDelete(ctx context.Context, id string) (int64, *apierrors.APIError) {
	var count int64
	errI := withretry.WithRetry(ctx, withretry.DefaultConfig, transientError, func(ctx context.Context) error {
		result, err := d.db.RefreshTokenDelete(ctx, id)
		count = result
		return err
	})
	return count, transformError(errI)
}

func (d *Psql) RefreshTokenDeleteFromUser(ctx context.Context, userID string) (int64, *apierrors.APIError) {
	var count int64
	errI := withretry.WithRetry(ctx, withretry.DefaultConfig, transientError, func(ctx context.Context) error {
		result, err := d.db.RefreshTokenDeleteFromUser(ctx, userID)
		count = result
		return err
	})
	return count, transformError(errI)
}
