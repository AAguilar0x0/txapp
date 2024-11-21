package bootstrap

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
	data, err := init(d.services)
	if err == nil {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		done := make(chan bool)

		defer func() {
			signal.Stop(sigCh)
			close(sigCh)
			close(done)
		}()

		go func() {
			data.Run()
			select {
			case done <- true:
			default:
			}
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
