package ksuid

import (
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/segmentio/ksuid"
)

type Ksuid struct{}

func New() (*Ksuid, error) {
	return &Ksuid{}, nil
}

func (d *Ksuid) Generate() (string, *apierrors.APIError) {
	id, err := ksuid.NewRandom()
	if err != nil {
		return "", apierrors.InternalServerError("cannot generate", err.Error())
	}
	return id.String(), nil
}
