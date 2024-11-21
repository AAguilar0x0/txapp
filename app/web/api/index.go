package api

import (
	"net/http"

	"github.com/AAguilar0x0/txapp/app/web/api/swagger"
	"github.com/AAguilar0x0/txapp/app/web/types"
	"github.com/AAguilar0x0/txapp/core/constants/envmodes"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/labstack/echo/v4"
)

type handler types.Handler

func New(e *echo.Group, h *types.Handler) {
	d := (*handler)(h)
	if h.Env == string(envmodes.Local) || h.Env == string(envmodes.Debug) {
		swagger.New(e.Group("/swagger"), h)
	}
	e.Use(
		h.Middlewares.RequestLogger(),
	)

	auth := e.Group("/auth")
	auth.POST("/signin", d.Signin)
}

type postSignin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// @Scheme http
// @Tags auth
// @Accept json
// @Produce json
// @Param body body postSignin true "Body"
// @Success 200 {object} string "Kani"
// @Router /auth/signin [post]
func (d *handler) Signin(c echo.Context) error {
	var body postSignin
	if err := c.Bind(&body); err != nil {
		return apierrors.BadRequest("Validation error")
	}
	if err := d.Vldtr.Struct(&body); err != nil {
		return err
	}
	aToken, rToken, err := d.Auth.SignIn(c.Request().Context(), body.Email, body.Password)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, map[string]any{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
	return nil
}
