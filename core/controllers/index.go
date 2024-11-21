package controllers

import (
	"github.com/AAguilar0x0/txapp/core/controllers/auth"
	"github.com/AAguilar0x0/txapp/core/services"
)

type DefaultControllerFactory struct {
	services services.ServiceProvider
}

func New(services services.ServiceProvider) *DefaultControllerFactory {
	d := DefaultControllerFactory{services}
	return &d
}

func (d *DefaultControllerFactory) Auth() (*auth.Auth, error) {
	db, err := d.services.Database()
	if err != nil {
		return nil, err
	}
	jwt, err := d.services.JWTokenizer()
	if err != nil {
		return nil, err
	}
	hash, err := d.services.Hash()
	if err != nil {
		return nil, err
	}
	return auth.New(db, jwt, hash)
}
