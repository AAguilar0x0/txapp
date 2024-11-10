package app

import (
	"github.com/AAguilar0x0/txapp/core/services"
)

type Initializer func(env services.Environment, config func(configs ...AppCallback)) Lifecycle

type Lifecycle interface {
	Run()
	Close()
}
