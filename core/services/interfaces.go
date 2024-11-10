package services

import "github.com/AAguilar0x0/txapp/core/pkg/apierrors"

type EnumValidator interface {
	ValidateEnum() bool
}

type Validator interface {
	Struct(s interface{}) *apierrors.APIError
}

type Authenticator interface {
	Hash(input string) (string, *apierrors.APIError)
	CompareHash(input, hash string) bool
	VerifyJWT(token string) *apierrors.APIError
}

type Environment interface {
	MustPresent(key ...string)
	MustGet(key string) string
	Get(key string) string
	GetDefault(key string, defaultValue string) string
}

type IDGenerator interface {
	Generate() (string, *apierrors.APIError)
}
