package main

import (
	"context"
	"embed"
	"net/http"
	"sync"
	"time"

	"github.com/AAguilar0x0/txapp/bootstrap"
	"github.com/AAguilar0x0/txapp/cmd/web/api"
	"github.com/AAguilar0x0/txapp/cmd/web/constants"
	"github.com/AAguilar0x0/txapp/cmd/web/pages"
	"github.com/AAguilar0x0/txapp/cmd/web/pkg/cookiemngr"
	"github.com/AAguilar0x0/txapp/cmd/web/pkg/middlewares"
	"github.com/AAguilar0x0/txapp/cmd/web/pkg/vfs"
	"github.com/AAguilar0x0/txapp/cmd/web/types"
	"github.com/AAguilar0x0/txapp/core/controllers"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	srvc "github.com/AAguilar0x0/txapp/core/services"
	"github.com/AAguilar0x0/txapp/extern"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

//go:embed static/*
var static embed.FS

type Web struct {
	wg   *sync.WaitGroup
	e    *echo.Echo
	port string
}

func New(services bootstrap.ServiceProvider) (bootstrap.Lifecycle, error) {
	env, serr := services.Environment()
	if serr != nil {
		return nil, serr
	}
	validator, serr := services.Validator()
	if serr != nil {
		return nil, serr
	}
	enc, serr := services.Encryptor()
	if serr != nil {
		return nil, serr
	}

	ed25519, err := enc.Asymmetric(srvc.EncryptEd25519)
	if err != nil {
		return nil, err
	}
	jwt, err := ed25519.PrivateKey([]byte(env.MustGet("AUTH_SECRET")))
	if err != nil {
		return nil, err
	}

	controllers := controllers.New(services)
	user, serr := controllers.Auth(jwt)
	if serr != nil {
		return nil, serr
	}
	vfs, serr := vfs.New(static, constants.StaticRoute, constants.StaticFilesMaxAge)
	if serr != nil {
		return nil, serr
	}
	cookie := cookiemngr.NewCookieManager("", true)

	d := Web{
		wg:   &sync.WaitGroup{},
		port: env.GetDefault("PORT", "8080"),
	}
	h := &types.Handler{
		Env:         env.Get("ENV"),
		Wg:          d.wg,
		Vldtr:       validator,
		Auth:        user,
		Cookie:      cookie,
		Middlewares: middlewares.New(cookie, user, vfs),
		VFS:         vfs,
	}

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

	d.e.Use(
		h.Middlewares.RateLimit(),
		h.Middlewares.RemoveTrailingSlash(),
		h.Middlewares.BodyDump(h.Env),
	)

	// d.e.Static("/static", "cmd/web/static")
	// d.e.GET("/static/*", func(c echo.Context) error {
	// 	path := "cmd/web/static/" + c.Param("*")
	// 	etag := h.Static.FilesChecksum(c.Param("*"))
	// 	if match := c.Request().Header.Get("If-None-Match"); match != "" && match == etag {
	// 		return c.NoContent(http.StatusNotModified)
	// 	}
	// 	c.Response().Header().Set("Cache-Control", "public, max-age=31536000")
	// 	c.Response().Header().Set("ETag", etag)
	// 	return c.File(path)
	// })
	d.e.GET(constants.StaticRoute+"/*", echo.WrapHandler(h.VFS))

	pages.New(d.e.Group(""), h)
	api.New(d.e.Group("/api"), h)

	return &d, nil
}

func (d *Web) Run() {
	if err := d.e.Start(":" + d.port); err != nil && err != http.ErrServerClosed {
		d.e.Logger.Fatalj(log.JSON{
			"message": "shutting down the server",
			"error":   err.Error(),
		})
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
	a := bootstrap.New(extern.New())
	a.Start(New)
}
