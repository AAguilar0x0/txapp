package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/AAguilar0x0/txapp/app"
	"github.com/AAguilar0x0/txapp/cmd/web/api"
	"github.com/AAguilar0x0/txapp/cmd/web/pages"
	"github.com/AAguilar0x0/txapp/cmd/web/types"
	"github.com/AAguilar0x0/txapp/core/controllers/user"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Web struct {
	wg   *sync.WaitGroup
	e    *echo.Echo
	port string
}

func New(env services.Environment, config func(configs ...app.AppCallback)) app.Lifecycle {
	d := Web{
		wg: &sync.WaitGroup{},
	}
	h := &types.Handler{
		Env: env.Get("ENV"),
		Wg:  d.wg,
	}

	config(
		app.Validator(func(data services.Validator) {
			h.Vldtr = data
		}),
		app.UserController(func(data *user.User) {
			h.User = data
		}),
	)

	d.port = env.GetDefault("PORT", "8080")
	d.e = echo.New()

	d.e.HTTPErrorHandler = func(err error, c echo.Context) {
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
	d.e.Use(middleware.RemoveTrailingSlash())
	d.e.Static("/static", "cmd/web/static")
	pages.New(d.e.Group(""), h)
	api.New(d.e.Group("/api"), h)

	return &d
}

func (d *Web) Run() {
	if err := d.e.Start(":" + d.port); err != nil && err != http.ErrServerClosed {
		d.e.Logger.Fatal("shutting down the server")
	}
}

func (d *Web) Close() {
	d.wg.Wait()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := d.e.Shutdown(ctx); err != nil {
		d.e.Logger.Fatal(err)
	}
}

// @title WebApp
// @version 1.0
// @description This is the backend api for WebApp.

// @contact.name WebApp

// @host localhost:8080
// @BasePath /api
func main() {
	a := app.New()
	a.Start(New)
}
