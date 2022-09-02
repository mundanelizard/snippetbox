package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"snippetbox.mundanelizard.com/ui"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}

// When building your own version of alice use 'curry' or 'functor' or 'monad' or 'chain' or 'pipe'
// and the last function that should be call should be named 'fold' or 'reduce'
// then the function to add a new item to the 'curry' list should be called 'map'

// When building your own version of httprouter, include alice like feature in the application
// and the `.New` function should be able to accept default middlewares. It should also have
// the ability to create a subset, like in gin, but it should always return a native `http.Handler`

// After creating the httprouter and chain, create an MVP platform that automagically sets the
// directory structure and development environment for the user, and also has a cli tool.

// Create a package that works like the internal/validator package and the validator@v10 package,
// but it should allow you to be able to register extra functions and parsers apart from the default.
// then it would send back the error on in a map[string][]string. You can specify to terminate
// on the first instance of error when validator or label some fields as optional.

// Build your own session manager for your server
// Add csrf blockers check gorilla/csrf or justinas/nosurf
