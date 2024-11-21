package golangjwt

import (
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/AAguilar0x0/txapp/core/constants"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Jwt struct {
	pkey  crypto.PublicKey
	skey  crypto.PrivateKey
	idGen services.IDGenerator
}

func New(env services.Environment, idGen services.IDGenerator) (*Jwt, error) {
	appSecret := env.MustGet("AUTH_SECRET")
	d := Jwt{idGen: idGen}
	err := d.parseKeyPairEdPrivateKeyFromPEM([]byte(appSecret))
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (*Jwt) Close() error {
	return nil
}

func (d *Jwt) parseKeyPairEdPrivateKeyFromPEM(key []byte) *apierrors.APIError {
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return apierrors.InternalServerError("invalid key: Key must be a PEM encoded PKCS1 or PKCS8 key")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return apierrors.InternalServerError(err.Error())
	}

	var privateKey ed25519.PrivateKey
	var ok bool
	if privateKey, ok = parsedKey.(ed25519.PrivateKey); !ok {
		return apierrors.InternalServerError("key is not a valid Ed25519 private key")
	}

	d.skey = privateKey
	d.pkey = privateKey.Public()

	return nil
}

func (d *Jwt) GetJWTSubjectID(token string) (string, string, *apierrors.APIError) {
	t, err := d.verifyJWT(token, expiredToken)
	if err != nil {
		return "", "", err
	}
	return t.Subject, t.ID, nil
}

func (d *Jwt) GenerateToken(id, role string, durationMinutes uint, HS512Key string) (string, *apierrors.APIError) {
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(durationMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   id,
		},
	}
	var method jwt.SigningMethod = jwt.SigningMethodEdDSA
	if HS512Key != "" {
		method = jwt.SigningMethodHS512
	}
	token := jwt.NewWithClaims(method, claims)
	var k interface{}
	if HS512Key == "" {
		k = d.skey
	} else {
		k = []byte(HS512Key)
	}
	tokenStr, err := token.SignedString(k)
	if err != nil {
		return "", apierrors.InternalServerError("cannot generate", err.Error())
	}
	return tokenStr, nil
}

func (d *Jwt) verifyJWT(token string, errorFilters ...errorFilter) (*Claims, *apierrors.APIError) {
	claims := Claims{}
	jwt, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return d.pkey, nil
	})
	returnError := true
	if err != nil {
		for _, ev := range errorFilters {
			if ev(err) {
				returnError = false
				break
			}
		}
		if returnError {
			return &claims, apierrors.Unauthorized("Invalid token", err.Error())
		}
	}
	if !jwt.Valid && returnError {
		return &claims, apierrors.Unauthorized("Invalid token")
	}
	return &claims, nil
}

func (d *Jwt) VerifyJWT(token string) *apierrors.APIError {
	_, err := d.verifyJWT(token)
	return err
}

func (d *Jwt) IsAccessTokenValid(accessToken, refreshToken string) (*services.TokenValid, *apierrors.APIError) {
	data := services.TokenValid{}
	rToken, err := d.verifyJWT(refreshToken)
	if err != nil && expiredToken(err) {
		data.RefreshTokenID = rToken.ID
		return &data, err
	} else if err != nil {
		return nil, err
	}
	data.UserID = rToken.Subject
	data.RefreshTokenID = rToken.ID
	aToken, err := d.verifyJWT(accessToken, expiredToken)
	if err != nil {
		return &data, err
	}
	if time.Now().Before(aToken.ExpiresAt.Time) {
		return &data, nil
	}
	data.Expired = true
	if aToken.Subject != rToken.Subject {
		return &data, apierrors.Forbidden("User mismatch")
	}
	if aToken.Issuer != rToken.ID {
		return &data, apierrors.Forbidden("Tokens mismatch")
	}
	return &data, nil
}

func (d *Jwt) GenerateAuthTokens(id, role string) (*services.AuthTokens, *apierrors.APIError) {
	jwtID, err := d.idGen.Generate()
	if err != nil {
		return nil, err
	}
	var method jwt.SigningMethod = jwt.SigningMethodEdDSA
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.RTokenDaysDuration * time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   id,
			ID:        jwtID,
		},
	}
	refreshToken, errI := jwt.NewWithClaims(method, claims).SignedString(d.skey)
	if errI != nil {
		return nil, apierrors.InternalServerError("Cannot generate tokens", err.Error())
	}
	claims = Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.ATokenMinutesDuration * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   id,
			Issuer:    jwtID,
		},
	}
	accessToken, errI := jwt.NewWithClaims(method, claims).SignedString(d.skey)
	if errI != nil {
		return nil, apierrors.InternalServerError("Cannot generate tokens", err.Error())
	}
	return &services.AuthTokens{RefreshTokenID: jwtID, RefreshToken: refreshToken, RefreshTokenExpiresAt: claims.ExpiresAt.Time, AccessToken: accessToken}, nil
}
