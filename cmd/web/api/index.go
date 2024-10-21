package api

import (
	"net/http"

	"github.com/AAguilar0x0/bapp/cmd/web/types"
	"github.com/labstack/echo/v4"
)

func Setup(g *echo.Group, h *types.Handler) {
	auth := g.Group("/auth")
	auth.POST("/signin", func(c echo.Context) error {
		c.String(http.StatusOK, "Kani")
		// h.User.SignIn(c.Request().Context(),)
		return nil
	})
}
