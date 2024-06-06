package main

// Imported the required packages
import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

// Created a routes function whch will return the routes
func (app *application) routes() http.Handler {
	// Created a NewServeMux instance and Handle the functions below
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.Then(http.HandlerFunc(app.Home)))
	mux.Post("/addTask", dynamicMiddleware.Then(http.HandlerFunc(app.AddList)))
	mux.Post("/deleteTask", dynamicMiddleware.Then(http.HandlerFunc(app.DeleteList)))
	mux.Post("/updateTask", dynamicMiddleware.Then(http.HandlerFunc(app.updateList)))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))
	// Serving Files Which is used to make the styles of the site
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	// return app.logRequest(secureHeaders(app.logRespond(mux)))
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	return standardMiddleware.Then(mux)
}
