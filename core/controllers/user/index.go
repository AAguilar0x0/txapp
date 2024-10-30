package user

import (
	"context"

	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/AAguilar0x0/txapp/extern/db/psql"
)

type User struct {
	db   *psql.Queries
	auth services.Authenticator
}

func New(db *psql.Queries, auth services.Authenticator) (*User, error) {
	user := User{
		db,
		auth,
	}
	return &user, nil
}

func (d *User) SignIn(ctx context.Context, email, password string) error {
	user, err := d.db.GetUserForAuth(ctx, email)
	if err != nil {
		return apierrors.InternalServerError("Error getting user", err.Error())
	}
	if !d.auth.CompareHash(password, user.Password) {
		return apierrors.Unauthorized("Invalid password")
	}
	return nil
}
