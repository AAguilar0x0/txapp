package psql

import (
	"context"

	"github.com/AAguilar0x0/txapp/core/models"
)

func (d *Psql) CreateUser(ctx context.Context, email string, firstName string, lastName string, password string, role string) (models.User, error) {
	data, err := d.db.CreateUser(ctx, email, firstName, lastName, password, role)
	return models.User(data), transformError(err)
}

func (d *Psql) GetUserForAuth(ctx context.Context, email string) (models.GetUserForAuthRow, error) {
	data, err := d.db.GetUserForAuth(ctx, email)
	return models.GetUserForAuthRow(data), transformError(err)
}
