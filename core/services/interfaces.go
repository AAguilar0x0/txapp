package services

import (
	"io"
	"time"

	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
)

type ServiceProvider interface {
	Environment() (Environment, error)
	Database() (models.Database, error)
	Migrator() (models.Migrator, error)
	Validator() (Validator, error)
	JWTokenizer() (JWTokenizer, error)
	Hash() (Hash, error)
	IDGenerator() (IDGenerator, error)
	io.Closer
}

type EnumValidator interface {
	ValidateEnum() bool
}

type Validator interface {
	Struct(s interface{}) *apierrors.APIError
	Var(f interface{}, tag string) *apierrors.APIError
}

type AuthTokens struct {
	RefreshTokenID        string
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	AccessToken           string
}

type TokenValid struct {
	Expired        bool
	UserID         string
	RefreshTokenID string
}

type JWTokenizer interface {
	GetJWTSubjectID(token string) (string, string, *apierrors.APIError)
	GenerateToken(id, role string, durationMinutes uint, HS512Key string) (string, *apierrors.APIError)
	VerifyJWT(token string) *apierrors.APIError
	IsAccessTokenValid(accessToken, refreshToken string) (*TokenValid, *apierrors.APIError)
	GenerateAuthTokens(id, role string) (*AuthTokens, *apierrors.APIError)
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
