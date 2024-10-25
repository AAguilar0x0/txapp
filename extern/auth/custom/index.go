package authcustom

import (
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Auth struct {
	init bool
	pkey crypto.PublicKey
	skey crypto.PrivateKey
}

func (d *Auth) Init(env services.Environment) error {
	if d.init {
		return nil
	}
	appSecret := env.MustGet("AUTH_SECRET")
	err := d.parseKeyPairEdPrivateKeyFromPEM([]byte(appSecret))
	if err != nil {
		return err
	}
	d.init = true
	return nil
}

func (*Auth) Close() {}

func (d *Auth) Hash(input string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(input), 15)
	return string(bytes), err
}

func (d *Auth) CompareHash(input, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	return err == nil
}

func (d *Auth) parseKeyPairEdPrivateKeyFromPEM(key []byte) error {
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return errors.New("invalid key: Key must be a PEM encoded PKCS1 or PKCS8 key")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	var privateKey ed25519.PrivateKey
	var ok bool
	if privateKey, ok = parsedKey.(ed25519.PrivateKey); !ok {
		return errors.New("key is not a valid Ed25519 private key")
	}

	d.skey = privateKey
	d.pkey = privateKey.Public()

	return nil
}

func (d *Auth) GenerateToken(id, role, key string) (string, error) {
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
	return token.SignedString(k)
}

func (d *Auth) VerifyJWT(token string) error {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return d.pkey, nil
	})
	if err != nil {
		return err
	}
	if !t.Valid {
		return errors.New("invalid token")
	}
	return nil
}
