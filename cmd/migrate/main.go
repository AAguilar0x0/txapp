package main

import (
	"log/slog"
	"strconv"

	"github.com/AAguilar0x0/txapp/app"
	"github.com/AAguilar0x0/txapp/core/pkg/assert"
	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/AAguilar0x0/txapp/extern/db/psql"
)

type Migrate struct {
	dir     string
	command string
	version *int64
	db      *psql.Queries
}

func New(env services.Environment, config func(configs ...app.AppCallback)) app.Lifecycle {
	dir := env.GetDefault("dir", "cmd/migrate/migrations")
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
		dir:     dir,
		command: command,
		version: version,
	}
	config(
		app.Database(func(db *psql.Queries) {
			d.db = db
		}),
	)

	return &d
}

func (d *Migrate) Run() {
	if err := d.db.Migrate(d.dir, d.command, d.version, true); err != nil {
		slog.Error("Error running migration", "error", err.Error())
	}
}

func (d *Migrate) Close() {}

func main() {
	a := app.New()
	a.Start(New)
}
