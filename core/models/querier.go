package models

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, email string, firstName string, lastName string, password string, role string) (User, error)
	GetUserForAuth(ctx context.Context, email string) (GetUserForAuthRow, error)
}
