package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/onehappyfellow/daebak-web/context"
	"github.com/onehappyfellow/daebak-web/models"
	"github.com/onehappyfellow/daebak-web/views"
)

const sessionCookieName = "session"
const sessionCookieSecret = "replace-with-a-secret-key"

type UsersHtml struct {
	Templates struct {
		Register views.Template
		Login    views.Template
		Forgot   views.Template
		Reset    views.Template
	}
	UserService *models.UserService
}

func (c UsersHtml) Register(w http.ResponseWriter, r *http.Request) {
	var data struct{ Error string }
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		_, err := c.UserService.CreateUser(email, password)
		if err != nil {
			data.Error = "Registration failed"
		} else {
			setSessionCookie(w, email)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	c.Templates.Register.Execute(w, r, data)
}

func (c UsersHtml) Login(w http.ResponseWriter, r *http.Request) {
	var data struct{ Error string }
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		user, err := c.UserService.Authenticate(email, password)
		if err != nil {
			data.Error = "Invalid credentials"
		} else {
			setSessionCookie(w, user.Email)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	c.Templates.Login.Execute(w, r, data)
}

func (c UsersHtml) Logout(w http.ResponseWriter, r *http.Request) {
	clearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c UsersHtml) Forgot(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Sent bool
	}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		token, err := c.UserService.SetResetToken(email)
		if err == nil {
			// TODO: send email with reset link containing token
			fmt.Printf("Reset link for %s http://localhost:3000/users/reset?token=%s\n", email, token)
			data.Sent = true
		}
	}
	c.Templates.Forgot.Execute(w, r, data)
}

func (c UsersHtml) Reset(w http.ResponseWriter, r *http.Request) {
	var data struct{ Error string }
	token := r.URL.Query().Get("token")
	if r.Method == http.MethodPost {
		password := r.FormValue("password")
		user, err := c.UserService.GetByResetToken(token)
		if err != nil {
			data.Error = "Invalid or expired token"
		} else {
			err = c.UserService.ResetPassword(user.ID, password)
			if err != nil {
				data.Error = "Reset failed"
			} else {
				setSessionCookie(w, user.Email)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}
	}
	c.Templates.Reset.Execute(w, r, data)
}

func (u UsersHtml) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}

// --- Session helpers ---

func setSessionCookie(w http.ResponseWriter, email string) {
	sig := signSession(email)
	value := base64.StdEncoding.EncodeToString([]byte(email)) + "|" + sig
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})
}

func getSessionEmail(r *http.Request) string {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return ""
	}
	parts := strings.Split(cookie.Value, "|")
	if len(parts) != 2 {
		return ""
	}
	emailBytes, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return ""
	}
	email := string(emailBytes)
	if signSession(email) != parts[1] {
		return ""
	}
	return email
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
	})
}

func signSession(email string) string {
	h := hmac.New(sha256.New, []byte(sessionCookieSecret))
	h.Write([]byte(email))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

type UserMiddleware struct {
	UserService *models.UserService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := getSessionEmail(r)
		if email == "" {
			// No session cookie found, proceed without setting a user.
			next.ServeHTTP(w, r)
			return
		}
		// If we have an email, try to lookup the user with that email.
		user, err := umw.UserService.GetByEmail(email)
		if err != nil {
			// User not found or some other error, proceed without setting a user.
			next.ServeHTTP(w, r)
			return
		}
		// If we get to this point, we have a user that we can store in the context!
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/users/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
