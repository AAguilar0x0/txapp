package api

import (
	"net/http"

	"github.com/AAguilar0x0/bapp/cmd/web/types"
	"github.com/AAguilar0x0/bapp/pkg/apierrors"
	"github.com/labstack/echo/v4"
)

func Setup(g *echo.Group, h *types.Handler) {
	auth := g.Group("/auth")
	auth.POST("/signin", func(c echo.Context) error {
		type request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var body request
		if err := c.Bind(&body); err != nil {
			return apierrors.BadRequest("Validation error")
		}
		if err := h.Vldtr.Struct(&body); err != nil {
			return apierrors.BadRequest("Validation error")
		}
		if err := h.User.SignIn(c.Request().Context(), body.Email, body.Password); err != nil {
			return err
		}
		c.String(http.StatusOK, "Kani")
		return nil
	})
}
