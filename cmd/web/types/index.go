package types

import (
	"github.com/AAguilar0x0/bapp/core/controllers/user"
	"github.com/AAguilar0x0/bapp/core/services"
)

type Handler struct {
	Env   services.Environment
	Vldtr services.Validator
	User  *user.User
}
