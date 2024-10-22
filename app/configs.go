package app

import (
	"github.com/AAguilar0x0/bapp/core/controllers/user"
	"github.com/AAguilar0x0/bapp/core/pkg/assert"
	"github.com/AAguilar0x0/bapp/core/services"
	authcustom "github.com/AAguilar0x0/bapp/extern/auth/custom"
	"github.com/AAguilar0x0/bapp/extern/db/psql"
	"github.com/AAguilar0x0/bapp/extern/env"
	"github.com/AAguilar0x0/bapp/extern/validator/validatorv10"
)

// ==================================================================================== //
// SERVICES
// ==================================================================================== //

func environment(app *App) {
	app.Env = env.New()
}

func Database(cb func(db *psql.DB)) AppCallback {
	return func(app *App) {
		if temp := app.services.db; temp != nil {
			cb(temp)
			return
		}
		db, err := psql.New(
			app.Env.CommandLineFlagWithDefault("DB_HOST", "localhost"),
			app.Env.CommandLineFlagWithDefault("DB_USER", "postgres"),
			app.Env.CommandLineFlagWithDefault("DB_PASSWORD", "postgres"),
			app.Env.CommandLineFlagWithDefault("DB_NAME", "postgres"),
			app.Env.CommandLineFlagWithDefault("DB_PORT", "5432"),
			app.Env.CommandLineFlagWithDefault("DB_SSLMODE", "disable"),
		)
		assert.NoError(err, "psql instantiation", "fault", "New")
		cb(db)
		app.CleanUp(func(app *App) {
			db.Close()
		})
	}
}

func Auth(cb func(auth services.Authenticator)) AppCallback {
	return func(app *App) {
		if temp := app.services.auth; temp != nil {
			cb(temp)
			return
		}
		appSecret := app.Env.CommandLineFlagPanics("AUTH_SECRET")
		auth, err := authcustom.New([]byte(appSecret))
		assert.NoError(err, "authcustom instantiation", "fault", "New")
		cb(auth)
	}
}

func Validator(cb func(data services.Validator)) AppCallback {
	return func(app *App) {
		data := validatorv10.New()
		cb(data)
	}
}

// ==================================================================================== //
// CONTROLLERS
// ==================================================================================== //

func UserController(cb func(data *user.User)) AppCallback {
	return func(app *App) {
		if app.services.db == nil {
			Database(func(db *psql.DB) {
				app.services.db = db
			})(app)
		}
		if app.services.auth == nil {
			Auth(func(auth services.Authenticator) {
				app.services.auth = auth
			})(app)
		}
		temp, err := user.New(app.services.db, app.services.auth)
		assert.NoError(err, "user instantiation", "fault", "New")
		cb(temp)
	}
}
