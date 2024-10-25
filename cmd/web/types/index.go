package types

import (
	"sync"

	"github.com/AAguilar0x0/txapp/core/controllers/user"
	"github.com/AAguilar0x0/txapp/core/services"
)

type Handler struct {
	Env   string
	Wg    *sync.WaitGroup
	Vldtr services.Validator
	User  *user.User
}
