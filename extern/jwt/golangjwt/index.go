package golangjwt

import (
	"time"

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

func (d *Claims) GetRole() string {
	return d.Role
}

type JwtEncryption struct {
	method jwt.SigningMethod
	idGen  services.IDGenerator
}

func New(idGen services.IDGenerator) (*JwtEncryption, error) {
	return &JwtEncryption{idGen: idGen}, nil
}

func (d *JwtEncryption) Asymmetric(enc services.Encryption) (services.AsymEncryptor, *apierrors.APIError) {
	data := JwtEncryption{
		idGen: d.idGen,
	}
	switch enc {
	case services.EncryptEd25519:
		data.method = jwt.SigningMethodEdDSA
		return &encEd25519{&data}, nil
	case services.EncryptRSA:
		data.method = jwt.SigningMethodRS512
		return &encRSA{&data}, nil
	default:
		return nil, apierrors.NotImplemented("Encryption method is not yet implemented")
	}
}

func (d *JwtEncryption) Symmetric(enc services.Encryption) (services.SymEncryptor, *apierrors.APIError) {
	data := JwtEncryption{
		idGen: d.idGen,
	}
	switch enc {
	case services.EncryptHS512:
		data.method = jwt.SigningMethodHS512
	default:
		return nil, apierrors.NotImplemented("Encryption method is not yet implemented")
	}
	return &symEncryption{&data}, nil
}
