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
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}

	// middleware

	// controllers
	userController := controllers.Users{
		UserService:    userService,
		SessionService: sessionService,
	}
	userController.Templates.Signup = views.Must(views.ParseFS(
		templates.FS, "layout.gohtml", "signup.gohtml",
	))
	userController.Templates.Signin = views.Must(views.ParseFS(
		templates.FS, "layout.gohtml", "signin.gohtml",
	))

	// setup router
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)

	// Routes
	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "layout.gohtml", "home.gohtml")),
	))
	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "layout.gohtml", "contact.gohtml")),
	))
	r.Get("/signup", userController.Signup)
	r.Post("/signup", userController.HandleSignup)
	r.Get("/login", userController.Signin)
	r.Post("/login", userController.HandleSignin)
	r.Get("/logout", userController.Logout)
	r.Get("/users/me", userController.CurrentUser) // TODO replace this
	r.NotFound(notFoundHandler)

	fmt.Println("Starting server on port 3000")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println(err)
	}
}
