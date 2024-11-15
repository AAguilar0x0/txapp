package services

import (
	"io"

	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
)

type ServiceProvider interface {
	Environment() (Environment, error)
	Database() (models.Database, error)
	Migrator() (models.Migrator, error)
	Validator() (Validator, error)
	Authenticator() (Authenticator, error)
	Hash() (Hash, error)
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
	GenerateToken(id, role, HS512Key string) (string, *apierrors.APIError)
	VerifyJWT(token string) *apierrors.APIError
}

type Hash interface {
	Hash(input string) (string, *apierrors.APIError)
	CompareHash(input, hash string) bool
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
