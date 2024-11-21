package cookiemngr

import (
	"net/http"
	"time"

	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
	"github.com/labstack/echo/v4"
)

type CookieManager struct {
	domain string
	secure bool
}

func NewCookieManager(domain string, secure bool) *CookieManager {
	return &CookieManager{
		domain: domain,
		secure: secure,
	}
}

func (d *CookieManager) Set(c echo.Context, name, value string, expires *time.Duration) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   d.domain,
		Secure:   d.secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	if expires != nil {
		cookie.Expires = time.Now().Add(*expires)
	}
	c.SetCookie(cookie)
}

func (d *CookieManager) Get(c echo.Context, name string) (string, *apierrors.APIError) {
	val, err := c.Cookie(name)
	if err != nil {
		return "", apierrors.InternalServerError("Cannot get cookie", err.Error())
	}
	return val.Value, nil
}

func (d *CookieManager) Delete(c echo.Context, name string) {
	expires := -24 * time.Hour
	d.Set(c, name, "", &expires)
}
