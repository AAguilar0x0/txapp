package app

import (
	"github.com/AAguilar0x0/bapp/core/controllers/user"
	"github.com/AAguilar0x0/bapp/core/services"
	authcustom "github.com/AAguilar0x0/bapp/extern/auth/custom"
	"github.com/AAguilar0x0/bapp/extern/db/psql"
	"github.com/AAguilar0x0/bapp/extern/env"
)

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
			"localhost",
			"postgres",
			"postgres",
			"postgres",
			"5432",
		)
		if err != nil {
			panic(err.Error())
		}
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
		if err != nil {
			panic(err.Error())
		}
		cb(auth)
	}
}

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
		if err != nil {
			panic(err.Error())
		}
		cb(temp)
	}
}
