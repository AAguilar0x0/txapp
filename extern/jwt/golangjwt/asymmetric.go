package golangjwt

import (
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/golang-jwt/jwt/v5"
)

type encEd25519 struct {
	*JwtEncryption
}

func (d *encEd25519) PrivateKey(key []byte) (services.JWTokenizer, *apierrors.APIError) {
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, apierrors.InternalServerError("Invalid key: key must be a PEM encoded PKCS1 or PKCS8 key")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, apierrors.InternalServerError(err.Error())
	}

	privateKey, ok := parsedKey.(ed25519.PrivateKey)
	if !ok {
		return nil, apierrors.InternalServerError("Key is not a valid Ed25519 private key")
	}

	return &jwtEncryption{d.JwtEncryption, privateKey, privateKey.Public()}, nil
}

func (d *encEd25519) PublicKey(key []byte) (services.JWTokenizer, *apierrors.APIError) {
	pubKey, err := jwt.ParseEdPublicKeyFromPEM(key)
	if err != nil {
		return nil, apierrors.InternalServerError(err.Error())
	}
	return &jwtEncryption{d.JwtEncryption, nil, pubKey}, nil
}

type encRSA struct {
	*JwtEncryption
}

func (d *encRSA) PrivateKey(key []byte) (services.JWTokenizer, *apierrors.APIError) {
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, apierrors.InternalServerError("Invalid key: key must be a PEM encoded PKCS1 or PKCS8 key")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, apierrors.InternalServerError(err.Error())
	}

	privateKey, ok := parsedKey.(rsa.PrivateKey)
	if !ok {
		return nil, apierrors.InternalServerError("Key is not a valid Ed25519 private key")
	}

	return &jwtEncryption{d.JwtEncryption, privateKey, privateKey.Public()}, nil
}

func (d *encRSA) PublicKey(key []byte) (services.JWTokenizer, *apierrors.APIError) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		return nil, apierrors.InternalServerError(err.Error())
	}
	return &jwtEncryption{d.JwtEncryption, nil, pubKey}, nil
}
