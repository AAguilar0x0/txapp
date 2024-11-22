package main

import (
	"embed"
	"log/slog"
	"strconv"

	"github.com/AAguilar0x0/txapp/bootstrap"
	"github.com/AAguilar0x0/txapp/core/pkg/assert"
	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/AAguilar0x0/txapp/extern"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Migrate struct {
	command string
	version *int64
	migrate services.Migrator
}

func New(services bootstrap.ServiceProvider) (bootstrap.Lifecycle, error) {
	env, err := services.Environment()
	if err != nil {
		return nil, err
	}
	migrate, err := services.Migrator()
	if err != nil {
		return nil, err
	}

	command := env.GetDefault("command", "up")
	versionStr := env.Get("version")
	var version *int64 = nil
	if versionStr != "" {
		ver, err := strconv.Atoi(versionStr)
		if err != nil {
			assert.NoError(err, "Invalid version number", "version", versionStr)
		}
		temp := int64(ver)
		version = &temp
	}
	d := Migrate{
		command: command,
		version: version,
		migrate: migrate,
	}

	return &d, nil
}

func (d *Migrate) Run() {
	if err := d.migrate.Migrate(migrations, "migrations", d.command, d.version, true); err != nil {
		slog.Error("Error running migration", "error", err.Error())
	}
}

func (d *Migrate) Close() {}

func main() {
	a := bootstrap.New(extern.New())
	a.Start(New)
}
