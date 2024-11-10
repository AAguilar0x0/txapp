package services

type EnumValidator interface {
	ValidateEnum() bool
}

type Validator interface {
	Struct(s interface{}) error
}

type Authenticator interface {
	Hash(input string) (string, error)
	CompareHash(input, hash string) bool
	VerifyJWT(token string) error
}

type Environment interface {
	MustPresent(key ...string)
	MustGet(key string) string
	Get(key string) string
	GetDefault(key string, defaultValue string) string
}

type IDGenerator interface {
	Generate() (string, error)
}
