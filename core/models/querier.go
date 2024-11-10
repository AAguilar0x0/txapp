package models

import (
	"context"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
)

type Querier interface {
	CreateUser(ctx context.Context, email string, firstName string, lastName string, password string, role string) (User, *apierrors.APIError)
	GetUserForAuth(ctx context.Context, email string) (GetUserForAuthRow, *apierrors.APIError)
}
