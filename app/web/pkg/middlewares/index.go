package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/AAguilar0x0/txapp/app/web/constants"
	"github.com/AAguilar0x0/txapp/app/web/pkg/cookiemngr"
	"github.com/AAguilar0x0/txapp/app/web/pkg/vfs"
	"github.com/AAguilar0x0/txapp/core/constants/envmodes"
	"github.com/AAguilar0x0/txapp/core/controllers/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type Middlewares struct {
	cookie *cookiemngr.CookieManager
	auth   *auth.Auth
	vfs    *vfs.VersionedFileServer
}

func New(cookie *cookiemngr.CookieManager, auth *auth.Auth, vfs *vfs.VersionedFileServer) *Middlewares {
	return &Middlewares{cookie, auth, vfs}
}

func (d *Middlewares) BodyDump(env string) echo.MiddlewareFunc {
	if env != envmodes.Debug {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}
	return middleware.BodyDump(func(ctx echo.Context, b1, b2 []byte) {
		slog.Info(ctx.Request().URL.Path+":BodyDump", "body", string(b1), "response", string(b2))
	})
}

func (d *Middlewares) RateLimit() echo.MiddlewareFunc {
	return middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20)))
}

func (d *Middlewares) RemoveTrailingSlash() echo.MiddlewareFunc {
	return middleware.RemoveTrailingSlash()
}

func (d *Middlewares) CSRF() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var stubNext echo.HandlerFunc = func(c echo.Context) error {
				return nil
			}
			mcsrf := middleware.CSRFWithConfig(middleware.CSRFConfig{
				TokenLookup:    "form:_csrf",
				CookiePath:     "/",
				CookieDomain:   "",
				CookieSecure:   true,
				CookieHTTPOnly: true,
				CookieSameSite: http.SameSiteStrictMode,
			})(stubNext)
			if err := mcsrf(c); err != nil {
				return err
			}
			val := c.Get(string(constants.CtxKeyCSRF))
			if csrf, ok := val.(string); ok {
				req := c.Request()
				newCtx := context.WithValue(req.Context(), constants.CtxKeyCSRF, csrf)
				c.SetRequest(req.WithContext(newCtx))
			}
			return next(c)
		}
	}
}

func (d *Middlewares) RequestLogger() echo.MiddlewareFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
	})
}

func (d *Middlewares) StaticFiles() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			newCtx := context.WithValue(req.Context(), constants.CtxKeyCSS, d.vfs.URLWithVersion("output.css"))
			newCtx = context.WithValue(newCtx, constants.CtxKeyHTMX, d.vfs.URLWithVersion("htmx.min.js"))
			c.SetRequest(req.WithContext(newCtx))
			return next(c)
		}
	}
}
