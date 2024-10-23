package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/AAguilar0x0/bapp/cmd/web/api/swagger"
	"github.com/AAguilar0x0/bapp/cmd/web/types"
	"github.com/AAguilar0x0/bapp/core/constants/envmodes"
	"github.com/AAguilar0x0/bapp/core/pkg/apierrors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type handler struct {
	*types.Handler
}

func New(e *echo.Group, h *types.Handler) {
	d := handler{Handler: h}
	if env := h.Env.CommandLineFlag("ENV"); env == string(envmodes.Local) || env == string(envmodes.Debug) {
		swagger.New(e.Group("/swagger"), h)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: false,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	auth := e.Group("/auth")
	auth.POST("/signin", d.Signin)
}

type postSignin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Scheme http
// @Tags auth
// @Accept json
// @Produce json
// @Param body body postSignin true "Body"
// @Success 200 {object} string "Kani"
// @Router /auth/signin [post]
func (h *handler) Signin(c echo.Context) error {
	var body postSignin
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
}
