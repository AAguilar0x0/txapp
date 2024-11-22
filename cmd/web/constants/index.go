package constants

type CtxKey string

const (
	CtxKeyCSRF CtxKey = "csrf"
	CtxKeyCSS  CtxKey = "css"
	CtxKeyHTMX CtxKey = "htmx"
)

const (
	StaticFilesMaxAge = 31536000
	StaticRoute       = "/static"
)
