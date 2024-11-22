package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/AAguilar0x0/txapp/core/constants/roles"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/services"
)

type Auth struct {
	db   services.Database
	jwt  services.JWTokenizer
	hash services.Hash
}

func New(db services.Database, jwt services.JWTokenizer, hash services.Hash) (*Auth, error) {
	auth := Auth{
		db,
		jwt,
		hash,
	}
	return &auth, nil
}

func (d *Auth) RefreshAuth(ctx context.Context, accessToken, refreshToken string) (string, string, *apierrors.APIError) {
	rToken, err := d.jwt.VerifyJWT(refreshToken)
	if err != nil {
		if sub := rToken.GetSub(); sub != "" {
			if _, err := d.db.RefreshTokenDeleteFromUser(ctx, sub); err != nil {
				return "", "", err
			}
		}
		return "", "", err
	}
	aToken, err := d.jwt.VerifyJWT(accessToken)
	if err != nil && err.Status != http.StatusForbidden {
		if _, err := d.db.RefreshTokenDeleteFromUser(ctx, rToken.GetSub()); err != nil {
			return "", "", err
		}
		return "", "", err
	}
	if time.Now().Before(aToken.GetExp()) {
		return accessToken, refreshToken, nil
	}
	if aToken.GetSub() != rToken.GetSub() {
		if _, err := d.db.RefreshTokenDeleteFromUser(ctx, rToken.GetSub()); err != nil {
			return "", "", err
		}
		return "", "", apierrors.Forbidden("User mismatch")
	}
	if aToken.GetIss() != rToken.GetID() {
		if _, err := d.db.RefreshTokenDeleteFromUser(ctx, rToken.GetSub()); err != nil {
			return "", "", err
		}
		return "", "", apierrors.Forbidden("Tokens mismatch")
	}

	if _, err = d.db.RefreshTokenGet(ctx, rToken.GetID(), rToken.GetSub()); err != nil && err.Status == http.StatusNotFound {
		if _, err := d.db.RefreshTokenDeleteFromUser(ctx, rToken.GetSub()); err != nil {
			return "", "", err
		}
		return "", "", apierrors.Forbidden("Use of invalidated token")
	} else if err != nil {
		return "", "", err
	}
	if _, err := d.db.RefreshTokenDelete(ctx, rToken.GetID()); err != nil {
		return "", "", err
	}
	role, err := d.db.UserGetForAuthID(ctx, rToken.GetSub())
	if err != nil {
		return "", "", err
	}
	tokens, err := d.jwt.GenerateAuthTokens(rToken.GetSub(), role)
	if err != nil {
		return "", "", err
	}
	if err := d.db.RefreshTokenCreate(ctx, tokens.RefreshToken.GetID(), tokens.RefreshToken.GetSub(), tokens.RefreshToken.GetExp()); err != nil {
		return "", "", err
	}
	return tokens.AccessJWT, tokens.RefreshJWT, nil
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
	if !d.hash.CompareHash(password, user.GetPassword()) {
		return "", "", apierrors.Unauthorized("Invalid password")
	}
	tokens, err := d.jwt.GenerateAuthTokens(user.GetID(), user.GetRole())
	if err != nil {
		return "", "", err
	}
	if err := d.db.RefreshTokenCreate(ctx, tokens.RefreshToken.GetID(), user.GetID(), tokens.RefreshToken.GetExp()); err != nil {
		return "", "", err
	}
	return tokens.AccessJWT, tokens.RefreshJWT, nil
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
