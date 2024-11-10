package ksuid

import (
	"github.com/segmentio/ksuid"
)

type Ksuid struct{}

func New() (*Ksuid, error) {
	return &Ksuid{}, nil
}

func (d *Ksuid) Generate() (string, error) {
	id, err := ksuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
