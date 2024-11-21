package middlewares

import (
	"net/http"
	"time"

	"github.com/AAguilar0x0/txapp/core/constants"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/labstack/echo/v4"
)

func (d *Middlewares) AuthGuest(c echo.Context, accessToken, refreshToken string) (string, *apierrors.APIError) {
	if refreshToken != "" {
		return "/home", nil
	}
	return "", nil
}

func (d *Middlewares) Authed(c echo.Context, accessToken, refreshToken string) (string, *apierrors.APIError) {
	if accessToken == "" || refreshToken == "" {
		d.cookie.Delete(c, "refresh_token")
		return "/", nil
	}
	aToken, rToken, err := d.auth.RefreshAuth(c.Request().Context(), accessToken, refreshToken)
	if err != nil {
		d.cookie.Delete(c, "access_token")
		d.cookie.Delete(c, "refresh_token")
		return "/", err
	}
	if accessToken == aToken && refreshToken == rToken {
		return "", nil
	}
	var expires *time.Duration
	remember, err := d.cookie.Get(c, "remember")
	if err == nil && remember == "true" {
		temp := constants.RTokenDaysDuration * time.Hour * 24
		expires = &temp
	}
	d.cookie.Set(c, "access_token", aToken, expires)
	d.cookie.Set(c, "refresh_token", rToken, expires)
	return "", nil
}

func (d *Middlewares) Auth(specifier func(c echo.Context, accessToken, refreshToken string) (string, *apierrors.APIError)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			aToken, err := d.cookie.Get(c, "access_token")
			if err != nil {
				aToken = ""
			}
			rToken, err := d.cookie.Get(c, "refresh_token")
			if err != nil {
				rToken = ""
			}
			redirectPath, err := specifier(c, aToken, rToken)
			if redirectPath != "" {
				if err := c.Redirect(http.StatusSeeOther, redirectPath); err != nil {
					return apierrors.InternalServerError("Cannot redirect", err.Error())
				}
				if err != nil {
					return err
				}
				return apierrors.Unauthorized("Unauthorized access")
			}
			if err != nil {
				return err
			}
			return next(c)
		}
	}
}
