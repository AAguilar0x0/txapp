package app

import (
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/AAguilar0x0/txapp/core/pkg/assert"
	"github.com/AAguilar0x0/txapp/core/services"
)

type AppCallback func(app *App)

type App struct {
	env      services.Environment
	closable []io.Closer
}

func New() *App {
	d := App{}
	d.config(environment)
	assert.Assert(d.env != nil, "unexpected Env nil value", "fault", "App.Env")
	return &d
}

func (d *App) registerResource(res io.Closer) {
	d.closable = append(d.closable, res)
}

func (d *App) config(configs ...AppCallback) {
	for _, conf := range configs {
		conf(d)
	}
}

func (d *App) Start(data Lifecycle) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	data.Init(d.env, d.config)

	go func() {
		data.Run()
		done <- true
	}()
	select {
	case <-sigCh:
	case <-done:
	}

	data.Close()
	for _, c := range d.closable {
		err := c.Close()
		assert.NoError(err, "resource close", "fault", "Close")
	}
}
