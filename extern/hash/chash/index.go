package chash

import (
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"golang.org/x/crypto/bcrypt"
)

type CHash struct{}

func New() (*CHash, error) {
	return &CHash{}, nil
}

func (d *CHash) Hash(input string) (string, *apierrors.APIError) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(input), 15)
	return string(bytes), apierrors.InternalServerError("cannot generate", err.Error())
}

func (d *CHash) CompareHash(input, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	return err == nil
}
