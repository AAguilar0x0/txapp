package types

import (
	"sync"

	"github.com/AAguilar0x0/txapp/core/controllers/user"
	"github.com/AAguilar0x0/txapp/core/services"
)

type Handler struct {
	Wg    *sync.WaitGroup
	Env   services.Environment
	Vldtr services.Validator
	User  *user.User
}
