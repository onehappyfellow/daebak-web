package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

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
	articleService := &models.ArticleService{DB: db}
	userService := &models.UserService{DB: db}
	tokenService := &models.TokenService{DB: db}
	vocabularyService := &models.VocabularyService{DB: db}
	// Instantiate new services for grammar, tags, and vocabulary (not yet used)
	_ = &models.GrammarService{DB: db}
	_ = &models.TagService{DB: db}

	// Set up middleware
	umw := controllers.UserMiddleware{
		UserService:  userService,
		TokenService: tokenService,
	}

	// controllers
	articlesJson := controllers.ArticlesJson{
		ArticleService: articleService,
	}
	vocabularyJson := controllers.VocabularyJson{
		VocabularyService: vocabularyService,
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
	adminHtml := controllers.AdminHtml{
		ArticleService:    articleService,
		VocabularyService: vocabularyService,
	}
	adminHtml.Templates.Form = views.Must(views.ParseFS(
		templates.FS, "layout.gohtml", "article-form.gohtml",
	))
	usersHtml := controllers.UsersHtml{
		UserService:  userService,
		TokenService: tokenService,
	}
	usersHtml.Templates.Register = views.Must(views.ParseFS(
		templates.FS, "layout.gohtml", "user-register.gohtml",
	))
	usersHtml.Templates.Login = views.Must(views.ParseFS(
		templates.FS, "layout.gohtml", "user-login.gohtml",
	))
	usersHtml.Templates.Forgot = views.Must(views.ParseFS(
		templates.FS, "layout.gohtml", "user-forgot.gohtml",
	))
	usersHtml.Templates.Reset = views.Must(views.ParseFS(
		templates.FS, "layout.gohtml", "user-reset.gohtml",
	))
	usersHtml.Templates.Current = views.Must(views.ParseFS(
		templates.FS, "layout.gohtml", "user-current.gohtml",
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
	r.Use(umw.SetUser)
	r.NotFound(notFoundHandler)

	// Public routes
	r.Get("/", articlesHtml.Home)
	r.Get("/a/{slug}", articlesHtml.Single)
	r.Get("/contact", controllers.StaticHandler("contact.gohtml"))
	r.Handle("/images/*", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	r.Get("/users/register", usersHtml.Register)
	r.Post("/users/register", usersHtml.Register)
	r.Get("/users/login", usersHtml.Login)
	r.Post("/users/login", usersHtml.Login)
	r.Get("/users/logout", usersHtml.Logout)
	r.Get("/users/forgot", usersHtml.Forgot)
	r.Post("/users/forgot", usersHtml.Forgot)
	r.Get("/users/reset", usersHtml.Reset)
	r.Post("/users/reset", usersHtml.Reset)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersHtml.CurrentUser)
		r.Post("/tokens", usersHtml.CurrentUser)
		r.Post("/tokens/delete", usersHtml.DeleteToken)
	})

	// Restricted routes
	r.Mount("/api/articles", apiRoutes(articlesJson))
	r.Mount("/api/vocabulary", vocabularyApiRoutes(vocabularyJson))
	r.Mount("/admin", adminRoutes(adminHtml))

	fmt.Println("Starting server on port 3000")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println(err)
	}
}

func vocabularyApiRoutes(c controllers.VocabularyJson) http.Handler {
	r := chi.NewRouter()
	r.Get("/", c.List)
	r.Post("/", c.Create)
	r.Post("/get-or-create", c.GetOrCreate)
	r.Put("/{id}", c.Update)
	r.Delete("/{id}", c.Delete)
	return r
}

func adminRoutes(c controllers.AdminHtml) http.Handler {
	r := chi.NewRouter()
	// TODO restrict to admin users
	r.Get("/articles/new", c.NewArticleForm)
	r.Get("/articles/{id}", c.EditArticleForm)
	r.Post("/images/upload", imageUploadHandler)
	return r
}

func apiRoutes(c controllers.ArticlesJson) http.Handler {
	r := chi.NewRouter()
	r.Get("/", c.GetAllArticles)
	r.Post("/", c.CreateArticle)
	r.Get("/{id}", c.GetArticle)
	r.Put("/{id}", c.UpdateArticle)
	r.Delete("/{id}", c.DeleteArticle)
	return r
}

// imageUploadHandler handles image uploads for admin users
func imageUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Limit upload size to 10MB
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File too large or invalid form", http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Only allow certain file extensions (basic check)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	ext := ""
	if len(handler.Filename) > 0 {
		for i := len(handler.Filename) - 1; i >= 0; i-- {
			if handler.Filename[i] == '.' {
				ext = handler.Filename[i:]
				break
			}
		}
	}
	if !allowed[ext] {
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}
	// Save file to images directory
	dst, err := os.OpenFile("images/"+handler.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	defer dst.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Image uploaded successfully as %s", handler.Filename)
}
