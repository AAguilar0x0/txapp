package auth

import (
	"context"
	"net/http"

	"github.com/AAguilar0x0/txapp/core/constants/roles"
	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/services"
)

type Auth struct {
	db   models.Database
	jwt  services.JWTokenizer
	hash services.Hash
}

func New(db models.Database, jwt services.JWTokenizer, hash services.Hash) (*Auth, error) {
	auth := Auth{
		db,
		jwt,
		hash,
	}
	return &auth, nil
}

func (d *Auth) RefreshAuth(ctx context.Context, accessToken, refreshToken string) (string, string, *apierrors.APIError) {
	result, err := d.jwt.IsAccessTokenValid(accessToken, refreshToken)
	if err != nil && result != nil {
		if result.UserID != "" {
			if _, err := d.db.RefreshTokenDeleteFromUser(ctx, result.UserID); err != nil {
				return "", "", err
			}
		} else if result.RefreshTokenID != "" {
			if _, err := d.db.RefreshTokenDelete(ctx, result.RefreshTokenID); err != nil {
				return "", "", err
			}
		}
		return "", "", err
	} else if err != nil {
		return "", "", err
	}
	if !result.Expired {
		return accessToken, refreshToken, nil
	}
	if _, err = d.db.RefreshTokenGet(ctx, result.RefreshTokenID, result.UserID); err != nil && err.Status == http.StatusNotFound {
		return "", "", apierrors.Forbidden("Use of invalidated token")
	} else if err != nil {
		return "", "", err
	}
	if _, err := d.db.RefreshTokenDelete(ctx, result.RefreshTokenID); err != nil {
		return "", "", err
	}
	role, err := d.db.UserGetForAuthID(ctx, result.UserID)
	if err != nil {
		return "", "", err
	}
	tokens, err := d.jwt.GenerateAuthTokens(result.UserID, role)
	if err != nil {
		return "", "", err
	}
	if err := d.db.RefreshTokenCreate(ctx, tokens.RefreshTokenID, result.UserID, tokens.RefreshTokenExpiresAt); err != nil {
		return "", "", err
	}
	return tokens.AccessToken, tokens.RefreshToken, nil
}

func (d *Auth) SignUp(ctx context.Context, email, password, fname, lname string) *apierrors.APIError {
	password, err := d.hash.Hash(password)
	if err != nil {
		return err
	}
	_, err = d.db.UserCreate(ctx, email, fname, lname, password, string(roles.User))
	if err != nil {
		return err
	}
	return nil
}

func (d *Auth) SignIn(ctx context.Context, email, password string) (string, string, *apierrors.APIError) {
	user, err := d.db.UserGetForAuth(ctx, email)
	if err != nil {
		return "", "", err
	}
	if !d.hash.CompareHash(password, user.Password) {
		return "", "", apierrors.Unauthorized("Invalid password")
	}
	tokens, err := d.jwt.GenerateAuthTokens(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}
	if err := d.db.RefreshTokenCreate(ctx, tokens.RefreshTokenID, user.ID, tokens.RefreshTokenExpiresAt); err != nil {
		return "", "", err
	}
	return tokens.AccessToken, tokens.RefreshToken, nil
}

func (d *Auth) SignOut(ctx context.Context, refreshToken string) *apierrors.APIError {
	_, id, err := d.jwt.GetJWTSubjectID(refreshToken)
	if err != nil {
		return err
	}
	count, err := d.db.RefreshTokenDelete(ctx, id)
	if err != nil {
		return err
	}
	if count != 1 {
		return apierrors.NotFound("Token not found")
	}
	return nil
}
