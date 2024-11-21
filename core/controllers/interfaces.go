package controllers

import (
	"github.com/AAguilar0x0/txapp/core/controllers/auth"
)

type ControllerFactory interface {
	Auth() (*auth.Auth, error)
}
