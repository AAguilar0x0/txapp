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
	PanicIfMissingEnvKey(key ...string)
	CommandLineFlag(key string) string
	CommandLineFlagWithDefault(key string, defaultValue string) string
	CommandLineFlagPanics(key string) string
}
