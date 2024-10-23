package envmodes

type EnvMode string

const (
	Prod  EnvMode = "prod"
	Test  EnvMode = "test"
	Dev   EnvMode = "dev"
	Local EnvMode = "local"
	Debug EnvMode = "debug"
)
