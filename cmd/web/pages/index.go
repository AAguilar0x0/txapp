package pages

import (
	"github.com/AAguilar0x0/bapp/cmd/web/pages/utils"
	"github.com/AAguilar0x0/bapp/cmd/web/types"
	"github.com/labstack/echo/v4"
)

type handler types.Handler

func New(e *echo.Group, h *types.Handler) {
	d := (*handler)(h)
	e.GET("", d.get)
}

func (d *handler) get(c echo.Context) error {
	return utils.Render(c, guest())
}
