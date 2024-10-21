package app

import (
	"context"
	"os"
	"os/signal"

	"github.com/AAguilar0x0/bapp/core/services"
	"github.com/AAguilar0x0/bapp/extern/db/psql"
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
	return &d
}

func (d *App) Config(configs ...AppCallback) *App {
	for _, conf := range configs {
		conf(d)
	}
	return d
}

func (d *App) CleanUp(cleanups ...AppCallback) *App {
	d.cleanupCB = append(d.cleanupCB, cleanups...)
	return d
}

func (d *App) Run(cb func()) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		cb()
	}()
	<-ctx.Done()
	for _, cb := range d.cleanupCB {
		cb(d)
	}
}
