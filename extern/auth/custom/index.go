package authcustom

import (
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Auth struct {
	pkey crypto.PublicKey
	skey crypto.PrivateKey
}

func New(env services.Environment) (*Auth, error) {
	appSecret := env.MustGet("AUTH_SECRET")
	d := Auth{}
	err := d.parseKeyPairEdPrivateKeyFromPEM([]byte(appSecret))
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (*Auth) Close() error {
	return nil
}

func (d *Auth) Hash(input string) (string, *apierrors.APIError) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(input), 15)
	return string(bytes), apierrors.InternalServerError("cannot generate", err.Error())
}

func (d *Auth) CompareHash(input, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	return err == nil
}

func (d *Auth) parseKeyPairEdPrivateKeyFromPEM(key []byte) *apierrors.APIError {
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

func (d *Auth) GenerateToken(id, role, key string) (string, *apierrors.APIError) {
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 15 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   id,
		},
	}
	var method jwt.SigningMethod = jwt.SigningMethodEdDSA
	if key != "" {
		method = jwt.SigningMethodHS512
	}
	token := jwt.NewWithClaims(method, claims)
	var k interface{}
	if key == "" {
		k = d.skey
	} else {
		k = []byte(key)
	}
	tokenStr, err := token.SignedString(k)
	if err != nil {
		return "", apierrors.InternalServerError("cannot generate", err.Error())
	}
	return tokenStr, nil
}

func (d *Auth) VerifyJWT(token string) *apierrors.APIError {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return d.pkey, nil
	})
	if err != nil {
		return apierrors.InternalServerError("cannot verify", err.Error())
	}
	if !t.Valid {
		return apierrors.Unauthorized("invalid token")
	}
	return nil
}
