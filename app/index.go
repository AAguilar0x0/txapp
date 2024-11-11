package app

import (
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/AAguilar0x0/txapp/core/pkg/assert"
	"github.com/AAguilar0x0/txapp/core/services"
)

type ServiceProvider interface {
	services.ServiceProvider
	io.Closer
}

type App struct {
	services ServiceProvider
}

func New(services ServiceProvider) *App {
	return &App{services}
}

func (d *App) Start(init Initializer) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	data, err := init(d.services)
	if err == nil {
		go func() {
			data.Run()
			done <- true
		}()
		select {
		case <-sigCh:
		case <-done:
		}

		data.Close()
	}

	err = d.services.Close()
	assert.NoError(err, "resource close", "fault", "Close")
}
