package main

import (
	"context"
	"net/http"
	"time"

	"github.com/AAguilar0x0/bapp/app"
	"github.com/AAguilar0x0/bapp/cmd/web/api"
	"github.com/AAguilar0x0/bapp/cmd/web/types"
	"github.com/AAguilar0x0/bapp/core/controllers/user"
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
	e.Use(middleware.Logger())
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
	a.Config(app.UserController(func(data *user.User) {
		h.User = data
	}))

	a.Run(func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	})
}
