package user

import (
	"context"

	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/services"
)

type User struct {
	db    models.Database
	auth  services.Authenticator
	idGen services.IDGenerator
}

func New(db models.Database, auth services.Authenticator, idGen services.IDGenerator) (*User, error) {
	user := User{
		db,
		auth,
		idGen,
	}
	return &user, nil
}

func (d *User) SignIn(ctx context.Context, email, password string) *apierrors.APIError {
	user, err := d.db.GetUserForAuth(ctx, email)
	if err != nil {
		return err
	}
	if !d.auth.CompareHash(password, user.Password) {
		return apierrors.Unauthorized("Invalid password")
	}
	return nil
}
