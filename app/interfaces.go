package app

import (
	"github.com/AAguilar0x0/txapp/core/services"
)

type Lifecycle interface {
	Init(env services.Environment, config func(configs ...AppCallback))
	Run()
	Close()
}

type Resource interface {
	Init(env services.Environment) error
	Close()
}
