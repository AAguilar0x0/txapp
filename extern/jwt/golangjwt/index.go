package golangjwt

import (
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/AAguilar0x0/txapp/core/constants"
	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (d *Claims) GetID() string {
	return d.ID
}

func (d *Claims) GetSub() string {
	return d.Subject
}

func (d *Claims) GetIss() string {
	return d.Issuer
}

func (d *Claims) GetExp() time.Time {
	return d.ExpiresAt.Time
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
	_, t, err := d.verifyJWT(token, expiredToken)
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

func (d *Jwt) verifyJWT(token string, errorFilters ...errorFilter) (*jwt.Token, *Claims, *apierrors.APIError) {
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
			return jwt, &claims, apierrors.Unauthorized("Invalid token", err.Error())
		}
	}
	if !jwt.Valid && returnError {
		return jwt, &claims, apierrors.Unauthorized("Invalid token")
	}
	return jwt, &claims, nil
}

func (d *Jwt) VerifyJWT(token string) (models.Token, *apierrors.APIError) {
	jwt, data, err := d.verifyJWT(token, expiredToken)
	if !jwt.Valid {
		err = apierrors.Forbidden("Token expired")
	}
	return data, err
}

func (d *Jwt) GenerateAuthTokens(id, role string) (*services.AuthTokens, *apierrors.APIError) {
	jwtID, err := d.idGen.Generate()
	if err != nil {
		return nil, err
	}
	var method jwt.SigningMethod = jwt.SigningMethodEdDSA
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
	refreshJWT, errI := jwt.NewWithClaims(method, rToken).SignedString(d.skey)
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
	accessJWT, errI := jwt.NewWithClaims(method, aToken).SignedString(d.skey)
	if errI != nil {
		return nil, apierrors.InternalServerError("Cannot generate tokens", err.Error())
	}
	return &services.AuthTokens{RefreshToken: &rToken, RefreshJWT: refreshJWT, AccessJWT: accessJWT}, nil
}
