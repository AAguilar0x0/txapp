package psql

import (
	"context"

	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/pkg/withretry"
)

func (d *Psql) UserCreate(ctx context.Context, email string, firstName string, lastName string, password string, role string) (*models.User, *apierrors.APIError) {
	id, err := d.idGen.Generate()
	if err != nil {
		return nil, err
	}
	var data *models.User
	errI := withretry.WithRetry(ctx, withretry.DefaultConfig, transientError, func(ctx context.Context) error {
		result, err := d.db.UserCreate(ctx, email, firstName, lastName, password, role, id)
		data = (*models.User)(result)
		return err
	})
	return data, transformError(errI)
}

func (d *Psql) UserGetForAuth(ctx context.Context, email string) (*models.UserGetForAuthRow, *apierrors.APIError) {
	data, err := d.db.UserGetForAuth(ctx, email)
	return (*models.UserGetForAuthRow)(data), transformError(err)
}

func (d *Psql) UserGetForAuthID(ctx context.Context, id string) (string, *apierrors.APIError) {
	data, err := d.db.UserGetForAuthID(ctx, id)
	return data, transformError(err)
}
