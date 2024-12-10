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
	articlesJson := controllers.ArticlesJson{
		ArticleService: articleService,
	}
	articlesHtml := controllers.ArticlesHtml{
		ArticleService: articleService,
	}
	articlesHtml.Templates.Single = views.Must(views.ParseFS(
		templates.FS, "layout.gohtml", "article.gohtml",
	))
	articlesHtml.Templates.List = views.Must(views.ParseFS(
		templates.FS, "layout.gohtml", "article-list.gohtml",
	))
	articlesHtml.Templates.Form = views.Must(views.ParseFS(
		templates.FS, "layout.gohtml", "article-form.gohtml",
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
	r.Get("/a/{slug}", articlesHtml.Single)
	r.Get("/articles/new", articlesHtml.CreateArticle)
	r.Post("/articles", articlesHtml.CreateArticle)
	r.Get("/trending", articlesHtml.Trending)
	r.Mount("/api/articles", apiRoutes(articlesJson))

	fmt.Println("Starting server on port 3000")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println(err)
	}
}

func apiRoutes(c controllers.ArticlesJson) http.Handler {
	r := chi.NewRouter()
	r.Get("/", c.GetAllArticles)
	r.Post("/", c.CreateArticle)
	r.Get("/{id}", c.GetArticle)
	r.Get("/slug/{slug}", c.GetArticleBySlug)
	r.Put("/{id}", c.UpdateArticle)
	r.Delete("/{id}", c.DeleteArticle)
	return r
}
