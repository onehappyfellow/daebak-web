package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/onehappyfellow/daebak-web/controllers"
	"github.com/onehappyfellow/daebak-web/templates"
	"github.com/onehappyfellow/daebak-web/views"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, `Page not found`, http.StatusNotFound)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.NotFound(notFoundHandler)

	// Routes
	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "layout.gohtml", "home.gohtml")),
	))
	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "layout.gohtml", "contact.gohtml")),
	))

	fmt.Println("Starting server on port 3000")
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println(err)
	}
}
