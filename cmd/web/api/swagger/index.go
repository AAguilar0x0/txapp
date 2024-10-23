package swagger

import (
	_ "github.com/AAguilar0x0/bapp/cmd/web/docs"
	"github.com/AAguilar0x0/bapp/cmd/web/types"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

// @title WebApp
// @version 1.0
// @description This is the backend api for WebApp.
// @contact.name WebApp
// @host localhost:8080

func Setup(e *echo.Group, h *types.Handler) {
	e.GET("/*", echoSwagger.WrapHandler)
}
