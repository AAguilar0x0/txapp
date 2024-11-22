package pages

import (
	"github.com/AAguilar0x0/txapp/cmd/web/pages/components"
	"github.com/AAguilar0x0/txapp/cmd/web/pages/utils"
	"github.com/AAguilar0x0/txapp/cmd/web/types"
	"github.com/labstack/echo/v4"
)

type handler types.Handler

func New(e *echo.Group, h *types.Handler) {
	d := (*handler)(h)
	e.GET("", d.get)
	e.Any("*", d.fallback)
}

func (d *handler) get(c echo.Context) error {
	return utils.Render(c, guest())
}

func (d *handler) fallback(c echo.Context) error {
	return utils.Render(c, components.PageNotFound())
}
