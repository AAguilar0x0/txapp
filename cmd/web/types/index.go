package types

import (
	"github.com/AAguilar0x0/txapp/core/controllers/user"
	"github.com/AAguilar0x0/txapp/core/services"
)

type Handler struct {
	Env   services.Environment
	Vldtr services.Validator
	User  *user.User
}
