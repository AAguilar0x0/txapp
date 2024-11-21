package types

import (
	"sync"

	"github.com/AAguilar0x0/txapp/app/web/pkg/cookiemngr"
	"github.com/AAguilar0x0/txapp/app/web/pkg/middlewares"
	"github.com/AAguilar0x0/txapp/app/web/pkg/vfs"
	"github.com/AAguilar0x0/txapp/core/controllers/auth"
	"github.com/AAguilar0x0/txapp/core/services"
)

type Handler struct {
	Env         string
	Wg          *sync.WaitGroup
	Vldtr       services.Validator
	Auth        *auth.Auth
	Cookie      *cookiemngr.CookieManager
	Middlewares *middlewares.Middlewares
	VFS         *vfs.VersionedFileServer
}
