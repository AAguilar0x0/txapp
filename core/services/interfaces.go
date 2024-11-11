package services

import (
	"io"

	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
)

type ServiceProvider interface {
	Environment() (Environment, error)
	Database() (models.Database, error)
	Validator() (Validator, error)
	Authenticator() (Authenticator, error)
	IDGenerator() (IDGenerator, error)
	io.Closer
}

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
