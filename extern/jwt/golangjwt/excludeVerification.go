package golangjwt

type errorFilter func(error) bool

func expiredToken(err error) bool {
	return err.Error() == "token has invalid claims: token is expired"
}
