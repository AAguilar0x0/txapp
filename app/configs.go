package app

import (
	"github.com/AAguilar0x0/txapp/core/controllers/user"
	"github.com/AAguilar0x0/txapp/core/pkg/assert"
	"github.com/AAguilar0x0/txapp/core/services"
	authcustom "github.com/AAguilar0x0/txapp/extern/auth/custom"
	"github.com/AAguilar0x0/txapp/extern/db/psql"
	"github.com/AAguilar0x0/txapp/extern/env"
	"github.com/AAguilar0x0/txapp/extern/validator/validatorv10"
)

// ==================================================================================== //
// SERVICES
// ==================================================================================== //

func environment(app *App) {
	app.env = env.New()
}

func Database(cb func(db *psql.Queries)) AppCallback {
	return func(app *App) {
		db, err := psql.NewDB(app.env)
		assert.NoError(err, "psql instantiation", "fault", "registerResource")
		app.registerResource(db)
		cb(db)
	}
}

func Auth(cb func(auth services.Authenticator)) AppCallback {
	return func(app *App) {
		auth, err := authcustom.New(app.env)
		assert.NoError(err, "authcustom instantiation", "fault", "registerResource")
		app.registerResource(auth)
		cb(auth)
	}
}

func Validator(cb func(data services.Validator)) AppCallback {
	return func(app *App) {
		validator, err := validatorv10.New(app.env)
		assert.NoError(err, "validatorv10 instantiation", "fault", "registerResource")
		app.registerResource(validator)
		cb(validator)
	}
}

// ==================================================================================== //
// CONTROLLERS
// ==================================================================================== //

func UserController(cb func(data *user.User)) AppCallback {
	return func(app *App) {
		var lDB *psql.Queries
		var lAuth services.Authenticator
		app.config(
			Database(func(db *psql.Queries) {
				lDB = db
			}),
			Auth(func(auth services.Authenticator) {
				lAuth = auth
			}),
		)
		temp, err := user.New(lDB, lAuth)
		assert.NoError(err, "user instantiation", "fault", "New")
		cb(temp)
	}
}
