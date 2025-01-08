package golangjwt

import "github.com/AAguilar0x0/txapp/core/services"

type symEncryption struct {
	*JwtEncryption
}

func (d *symEncryption) Key(key []byte) services.JWTokenizer {
	return &jwtEncryption{d.JwtEncryption, key, key}
}
