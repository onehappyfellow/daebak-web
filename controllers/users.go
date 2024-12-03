package controllers

import (
	"fmt"
	"net/http"

	"github.com/onehappyfellow/daebak-web/context"
	"github.com/onehappyfellow/daebak-web/errors"
	"github.com/onehappyfellow/daebak-web/models"
)

type Users struct {
	Templates struct {
		Signup         Template
		Signin         Template
		Profile        Template
		ForgotPassword Template
		ResetPassword  Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
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
		if errors.Is(err, models.ErrEmailTaken) {
			err = errors.Public(err, "That email is already associated with an account.")
		}
		if errors.Is(err, models.ErrPasswordInsecure) {
			err = errors.Public(err, fmt.Sprintf("Your password must be at least %d characters long.", models.MinPasswordLength))
		}
		data := struct{ Email string }{Email: email}
		u.Templates.Signup.Execute(w, r, data, err)
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
		if errors.Is(err, models.ErrInvalidAuth) {
			err = errors.Public(err, "That email or password is incorrect.")
		}
		data := struct{ Email string }{Email: email}
		u.Templates.Signin.Execute(w, r, data, err)
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

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	u.Templates.ForgotPassword.Execute(w, r, nil)
}

func (u Users) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	reset, err := u.PasswordResetService.Create(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	link := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", reset.Token)
	fmt.Println(link)
	// TODO send token via email
	msg := fmt.Sprintf("A password reset token has been sent to %s.", email)
	// TODO don't do this as an error
	u.Templates.ResetPassword.Execute(w, r, nil, errors.Public(fmt.Errorf("success"), msg))
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	password := r.FormValue("password")
	user, err := u.PasswordResetService.Consume(token)
	if err != nil {
		if errors.Is(err, models.ErrTokenInvalid) {
			err = errors.Public(err, "The password reset token is invalid.")
		}
		if errors.Is(err, models.ErrTokenExpired) {
			err = errors.Public(err, "The password reset token is expired.")
		}
		u.Templates.ResetPassword.Execute(w, r, nil, err)
		return
	}
	err = u.UserService.UpdatePassword(user, password)
	if err != nil {
		// bug: by this point consume succeeded so the token is no longer valid
		if errors.Is(err, models.ErrPasswordInsecure) {
			err = errors.Public(err, fmt.Sprintf("Your password must be at least %d characters long.", models.MinPasswordLength))
		}

		data := struct{ Token string }{Token: token}
		u.Templates.ResetPassword.Execute(w, r, data, err)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	SetCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/profile", http.StatusFound)
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
			// TOAST ERROR You must be logged in to access that page
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
