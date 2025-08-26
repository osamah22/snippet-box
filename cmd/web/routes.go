package main

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/osamah22/snippetbox/ui"

	"github.com/go-chi/chi"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	router.Use(app.recoverPanic)
	router.Use(app.logRequest)
	router.Use(secureHeaders)

	dynamic := alice.New(app.sessionManager.LoadAndSave).Append(app.authenticate, noSurf)
	protected := dynamic.Append(app.requireAuthentication)

	router.Get("/ping", ping)

	staticFS, err := fs.Sub(ui.Files, "static")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(staticFS))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	router.Handle("/static/", http.StripPrefix("/static", fileServer))
	router.Get("/", dynamic.ThenFunc(app.home).ServeHTTP)

	router.Get("/snippets/create", protected.ThenFunc(app.snippetCreaetView).ServeHTTP)
	router.Post("/snippets/create", protected.ThenFunc(app.snippetCreaet).ServeHTTP)
	router.Get("/snippets/view/{id}", dynamic.ThenFunc(app.snippetView).ServeHTTP)

	router.Get("/user/signup", dynamic.ThenFunc(app.userSignup).ServeHTTP)
	router.Post("/user/signup", dynamic.ThenFunc(app.userSignupPost).ServeHTTP)
	router.Get("/user/login", dynamic.ThenFunc(app.userLogin).ServeHTTP)
	router.Post("/user/login", dynamic.ThenFunc(app.userLoginPost).ServeHTTP)
	router.Post("/user/logout", protected.ThenFunc(app.userLogoutPost).ServeHTTP)

	return router
}
