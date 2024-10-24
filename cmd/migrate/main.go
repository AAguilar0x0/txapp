package main

import (
	"log/slog"
	"strconv"

	"github.com/AAguilar0x0/txapp/app"
	"github.com/AAguilar0x0/txapp/core/services"
	"github.com/AAguilar0x0/txapp/extern/db/psql"
)

func main() {
	a := app.New()
	var database *psql.DB
	a.Config(app.Database(func(db *psql.DB) {
		database = db
	}))
	a.Run(func(env services.Environment) {
		dir := env.CommandLineFlagWithDefault("dir", "cmd/migrate/migrations")
		command := env.CommandLineFlagWithDefault("command", "up")
		versionStr := env.CommandLineFlag("version")
		var version *int64 = nil
		if versionStr != "" {
			ver, err := strconv.Atoi(versionStr)
			if err != nil {
				slog.Error("Invalid version number", "version", version)
				return
			}
			temp := int64(ver)
			version = &temp
		}
		if err := database.Migrate(dir, command, version, true); err != nil {
			slog.Error("Error running migration", "error", err.Error())
		}
	})
}
