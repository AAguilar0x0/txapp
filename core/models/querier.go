package models

import (
	"context"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"time"
)

type Querier interface {
	RefreshTokenCreate(ctx context.Context, iD string, userID string, expiresAt time.Time) *apierrors.APIError
	RefreshTokenDelete(ctx context.Context, id string) (int64, *apierrors.APIError)
	RefreshTokenDeleteFromUser(ctx context.Context, userID string) (int64, *apierrors.APIError)
	RefreshTokenGet(ctx context.Context, iD string, userID string) (*RefreshToken, *apierrors.APIError)
	UserCreate(ctx context.Context, email string, firstName string, lastName string, password string, role string) (*User, *apierrors.APIError)
	UserGetForAuth(ctx context.Context, email string) (*UserGetForAuthRow, *apierrors.APIError)
	UserGetForAuthID(ctx context.Context, id string) (string, *apierrors.APIError)
}
