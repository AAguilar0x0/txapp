package services

import (
	"io"

	"github.com/AAguilar0x0/txapp/core/models"
	"github.com/AAguilar0x0/txapp/core/pkg/apierrors"
)

type ServiceProvider interface {
	Environment() (Environment, error)
	Database() (Database, error)
	Migrator() (Migrator, error)
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

// Encryption Start

type Encryption int

const (
	EncryptEd25519 Encryption = iota
	EncryptRSA
	EncryptHS512
)

type AuthTokens struct {
	RefreshToken models.Token
	RefreshJWT   string
	AccessJWT    string
}

type JWTokenizer interface {
	GetJWT(token string) (models.Token, *apierrors.APIError)
	GenerateToken(id, role string, durationMinutes uint) (string, *apierrors.APIError)
	VerifyJWT(token string) (models.Token, *apierrors.APIError)
	GenerateAuthTokens(id, role string) (*AuthTokens, *apierrors.APIError)
}

type AsymEncryptor interface {
	PrivateKey(key []byte) (JWTokenizer, *apierrors.APIError)
	PublicKey(key []byte) (JWTokenizer, *apierrors.APIError)
}

type SymEncryptor interface {
	Key(key []byte) JWTokenizer
}

type Encryptor interface {
	Asymmetric(Encryption) (AsymEncryptor, *apierrors.APIError)
	Symmetric(Encryption) (SymEncryptor, *apierrors.APIError)
}

// Encryption End

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
