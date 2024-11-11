package controllers

import (
	"github.com/AAguilar0x0/txapp/core/controllers/user"
	"github.com/AAguilar0x0/txapp/core/services"
)

type Controllers struct {
	services services.ServiceProvider
}

func New(services services.ServiceProvider) *Controllers {
	d := Controllers{services}
	return &d
}

func (d *Controllers) User() (*user.User, error) {
	db, err := d.services.Database()
	if err != nil {
		return nil, err
	}
	auth, err := d.services.Authenticator()
	if err != nil {
		return nil, err
	}
	return user.New(db, auth)
}
