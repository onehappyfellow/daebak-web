package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/onehappyfellow/daebak-web/controllers"
	"github.com/onehappyfellow/daebak-web/models"
	"github.com/onehappyfellow/daebak-web/templates"
	"github.com/onehappyfellow/daebak-web/views"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, `Page not found`, http.StatusNotFound)
}

func main() {
	// setup the database
	db, err := models.Open(models.DefaultPostresConfig())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// setup services
	articleService := &models.ArticleService{
		DB: db,
	}

	// controllers
	articleController := controllers.Articles{
		ArticleService: articleService,
	}
	articleController.Templates.Single = views.Must(views.ParseFS(
		templates.FS, "layout.gohtml", "article.gohtml",
	))
	articleController.Templates.List = views.Must(views.ParseFS(
		templates.FS, "layout.gohtml", "article-list.gohtml",
	))

	// setup router
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	// r.Use(middleware.URLFormat)
	// r.Use(middleware.Recoverer)
	// r.Use(middleware.Timeout(60 * time.Second))
	r.NotFound(notFoundHandler)

	// Routes
	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "layout.gohtml", "home.gohtml")),
	))
	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "layout.gohtml", "contact.gohtml")),
	))

	// r.With(articleController.SingleCtx).Get("/a/{articleSlug:[a-z-]+}", articleController.Single)
	// r.Get("/a/{articleSlug:[a-z-]+}", articleController.Single)

	r.Route("/a/{articleID}", func(r chi.Router) {
		r.Use(controllers.ArticleCtx)        // Load the *Article on the request context
		r.Get("/", articleController.Single) // GET /a/123
	})

	r.With(controllers.ArticleCtx).Get("/a/{articleSlug:[a-z-]+}", articleController.Single)
	// GET /a/whats-up

	r.Get("/trending", articleController.Trending)

	fmt.Println("Starting server on port 3000")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println(err)
	}
}
