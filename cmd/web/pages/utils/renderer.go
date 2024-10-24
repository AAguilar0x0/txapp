package utils

import (
	"github.com/AAguilar0x0/txapp/cmd/web/pages/components"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, comp templ.Component) error {
	return components.Main(comp).Render(c.Request().Context(), c.Response())
}
