package controllers

import (
	"fmt"
	"net/http"

	"github.com/onehappyfellow/daebak-web/context"
	"github.com/onehappyfellow/daebak-web/models"
)

type Users struct {
	Templates struct {
		Signup  Template
		Signin  Template
		Profile Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) Signup(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.Signup.Execute(w, r, data)
}

func (u Users) HandleSignup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Printf("could not create session for new user: %v\n", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	SetCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (u Users) Signin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.Signin.Execute(w, r, data)
}

func (u Users) HandleSignin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Printf("could not create session for authorized user: %v\n", err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	SetCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (u Users) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(CookieSession)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(c.Value)
	if err != nil {
		fmt.Printf("could not delete session %v\n", err)
	}
	DeleteCookie(w, CookieSession)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (u Users) Profile(w http.ResponseWriter, r *http.Request) {
	var data struct {
		User *models.User
	}
	data.User = context.User(r.Context())
	u.Templates.Profile.Execute(w, r, data)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (mw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := mw.SessionService.User(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (mw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
