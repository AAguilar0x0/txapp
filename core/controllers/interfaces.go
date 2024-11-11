package controllers

import (
	"github.com/AAguilar0x0/txapp/core/controllers/user"
)

type ControllerFactory interface {
	User() (*user.User, error)
}
