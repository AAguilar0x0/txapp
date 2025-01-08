package golangjwt

import (
	"crypto"
	"time"

	"github.com/AAguilar0x0/txapp/core/constants"
	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/golang-jwt/jwt/v5"
)

type jwtEncryption struct {
	*JwtEncryption
	privateKey crypto.PrivateKey
	publicKey  crypto.PublicKey
}

func (d *jwtEncryption) GetJWT(token string) (models.Token, *apierrors.APIError) {
	claims := Claims{}
	_, _, err := jwt.NewParser().ParseUnverified(token, &claims)
	if err != nil {
		return nil, apierrors.InternalServerError("Invalid token", err.Error())
	}
	return &claims, nil
}

func (d *jwtEncryption) GenerateToken(id, role string, durationMinutes uint) (string, *apierrors.APIError) {
	if d.privateKey == nil {
		return "", apierrors.NotImplemented("Method not implemented")
	}
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(durationMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   id,
		},
	}
	token := jwt.NewWithClaims(d.method, claims)
	tokenStr, err := token.SignedString(d.privateKey)
	if err != nil {
		return "", apierrors.InternalServerError("Cannot generate", err.Error())
	}
	return tokenStr, nil
}

func (d *jwtEncryption) VerifyJWT(token string) (models.Token, *apierrors.APIError) {
	if d.publicKey == nil {
		return nil, apierrors.NotImplemented("Method not implemented")
	}
	claims := Claims{}
	jwt, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return d.publicKey, nil
	})
	if err != nil {
		if err := err.Error(); err == "token has invalid claims: token is expired" {
			return &claims, apierrors.Forbidden("Token expired", err)
		}
		return &claims, apierrors.Unauthorized("Invalid token", err.Error())
	}
	if !jwt.Valid {
		return &claims, apierrors.Unauthorized("Invalid token")
	}
	return &claims, nil
}

func (d *jwtEncryption) GenerateAuthTokens(id, role string) (*services.AuthTokens, *apierrors.APIError) {
	if d.privateKey == nil {
		return nil, apierrors.NotImplemented("Method not implemented")
	}
	jwtID, err := d.idGen.Generate()
	if err != nil {
		return nil, err
	}
	rToken := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.RTokenDaysDuration * time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   id,
			ID:        jwtID,
		},
	}
	refreshJWT, errI := jwt.NewWithClaims(d.method, rToken).SignedString(d.privateKey)
	if errI != nil {
		return nil, apierrors.InternalServerError("Cannot generate tokens", err.Error())
	}
	aToken := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.ATokenMinutesDuration * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   id,
			Issuer:    jwtID,
		},
	}
	accessJWT, errI := jwt.NewWithClaims(d.method, aToken).SignedString(d.privateKey)
	if errI != nil {
		return nil, apierrors.InternalServerError("Cannot generate tokens", err.Error())
	}
	return &services.AuthTokens{RefreshToken: &rToken, RefreshJWT: refreshJWT, AccessJWT: accessJWT}, nil
}
