package psql

import (
	"context"

	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
)

func (d *Psql) CreateUser(ctx context.Context, id string, email string, firstName string, lastName string, password string, role string) (models.User, *apierrors.APIError) {
	data, err := d.db.CreateUser(ctx, id, email, firstName, lastName, password, role)
	return models.User(data), transformError(err)
}

func (d *Psql) GetUserForAuth(ctx context.Context, email string) (models.GetUserForAuthRow, *apierrors.APIError) {
	data, err := d.db.GetUserForAuth(ctx, email)
	return models.GetUserForAuthRow(data), transformError(err)
}
