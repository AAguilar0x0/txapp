package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/AAguilar0x0/bapp/core/services"
	"github.com/AAguilar0x0/bapp/extern/db/psql"
	"github.com/AAguilar0x0/bapp/pkg/assert"
)

type AppCallback func(app *App)

type App struct {
	Env      services.Environment
	services struct {
		db   *psql.DB
		auth services.Authenticator
	}
	cleanupCB []AppCallback
}

func New() *App {
	d := App{}
	d.Config(environment)
	assert.Assert(d.Env != nil, "unexpected Env nil value", "fault", "App.Env")
	return &d
}

func (d *App) Config(configs ...AppCallback) *App {
	for _, conf := range configs {
		conf(d)
	}
	return d
}

func (d *App) CleanUp(cleanups ...AppCallback) *App {
	assert.Assert(d.Env != nil, "unexpected Env nil value", "fault", "App.Env")
	d.cleanupCB = append(d.cleanupCB, cleanups...)
	return d
}

func (d *App) Run(cb func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		cb()
		done <- true
	}()
	select {
	case <-sigCh:
	case <-done:
	}
	for _, cb := range d.cleanupCB {
		cb(d)
	}
}
