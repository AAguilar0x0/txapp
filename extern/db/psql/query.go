package psql

import (
	"context"

	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
)

func (d *Psql) CreateUser(ctx context.Context, email string, firstName string, lastName string, password string, role string) (models.User, *apierrors.APIError) {
	id, err := d.idGen.Generate()
	if err != nil {
		return models.User{}, err
	}
	data, errI := d.db.CreateUser(ctx, email, firstName, lastName, password, role, id)
	return models.User(data), transformError(errI)
}

func (d *Psql) GetUserForAuth(ctx context.Context, email string) (models.GetUserForAuthRow, *apierrors.APIError) {
	data, err := d.db.GetUserForAuth(ctx, email)
	return models.GetUserForAuthRow(data), transformError(err)
}
