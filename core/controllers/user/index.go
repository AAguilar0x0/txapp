package user

import (
	"context"

	"github.com/AAguilar0x0/bapp/core/services"
	"github.com/AAguilar0x0/bapp/extern/db/psql"
	"github.com/AAguilar0x0/bapp/pkg/apierrors"
)

type User struct {
	db   *psql.DB
	auth services.Authenticator
}

func New(db *psql.DB, auth services.Authenticator) (*User, error) {
	user := User{
		db,
		auth,
	}
	return &user, nil
}

func (d *User) SignIn(ctx context.Context, email, password string) error {
	user, err := d.db.Instance().GetUserForAuth(ctx, email)
	if err != nil {
		return apierrors.InternalServerError("Error getting user", err.Error())
	}
	if !d.auth.CompareHash(password, user.Password) {
		return apierrors.Unauthorized("Invalid password")
	}
	return nil
}
