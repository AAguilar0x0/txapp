package swagger

import (
	_ "github.com/AAguilar0x0/txapp/app/web/docs"
	"github.com/AAguilar0x0/txapp/app/web/types"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func New(e *echo.Group, h *types.Handler) {
	e.GET("/*", echoSwagger.WrapHandler)
}
