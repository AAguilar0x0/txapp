package services

import (
	"context"
	"time"

	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
)

type Querier interface {
	// User
	UserCreate(ctx context.Context, email string, firstName string, lastName string, password string, role string) (models.User, *apierrors.APIError)
	UserGetForAuth(ctx context.Context, email string) (models.User, *apierrors.APIError)
	UserGetForAuthID(ctx context.Context, id string) (string, *apierrors.APIError)

	// RefreshToken
	RefreshTokenCreate(ctx context.Context, iD string, userID string, expiresAt time.Time) *apierrors.APIError
	RefreshTokenGet(ctx context.Context, iD string, userID string) (models.Token, *apierrors.APIError)
	RefreshTokenDelete(ctx context.Context, id string) (int64, *apierrors.APIError)
	RefreshTokenDeleteFromUser(ctx context.Context, userID string) (int64, *apierrors.APIError)
}
