package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AAguilar0x0/bapp/app"
	"github.com/AAguilar0x0/bapp/cmd/web/api"
	"github.com/AAguilar0x0/bapp/cmd/web/types"
	"github.com/AAguilar0x0/bapp/core/controllers/user"
	"github.com/AAguilar0x0/bapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/bapp/core/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	a := app.New()
	h := types.Handler{
		Env: a.Env,
	}

	port := a.Env.CommandLineFlagWithDefault("PORT", "8080")
	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		msg := http.StatusText(code)
		switch v := err.(type) {
		case *echo.HTTPError:
			code = v.Code
		case *apierrors.APIError:
			code = v.Status
			msg = v.Message
		}
		c.String(code, msg)
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
	e.Use(middleware.RemoveTrailingSlash())
	e.Static("/static", "cmd/web/static")
	api.Setup(e.Group("/api/v1"), &h)

	a.CleanUp(func(app *app.App) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	})
	a.Config(app.Validator(func(data services.Validator) {
		h.Vldtr = data
	}))
	a.Config(app.UserController(func(data *user.User) {
		h.User = data
	}))

	a.Run(func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	})
}
