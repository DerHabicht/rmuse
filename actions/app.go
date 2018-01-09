package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/envy"

	"github.com/derhabicht/rmuse/models"
	"github.com/gobuffalo/x/sessions"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			SessionName:  "_rmuse_session",
		})

		// Set the request content type to JSON
		app.Use(middleware.SetContentType("application/json"))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))

		app.GET("/", HomeHandler)

		// API V1 Grouping
		v1 := app.Group("/api/1")

		// Add middleware
		v1.Use(VerifyToken)
		v1.Middleware.Skip(VerifyToken, AuthCreateSession, UserCreate)

		// Login
		v1.POST("/login", AuthCreateSession)

		// Users
		v1.GET("/user", UserRead)
		v1.POST("/user", UserCreate)
	}

	return app
}
