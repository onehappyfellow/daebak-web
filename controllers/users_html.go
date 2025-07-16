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
		Current  views.Template
	}
	UserService  *models.UserService
	TokenService *models.TokenService
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

	var tokens []models.Token
	var err error
	if u.TokenService != nil {
		tokens, err = u.TokenService.ListByUserID(user.ID)
		if err != nil {
			tokens = nil
		}
	}

	var data struct {
		User   *models.User
		Tokens []models.Token
		Error  string
	}

	data.User = user
	data.Tokens = tokens

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		if name != "" && u.TokenService != nil {
			_, err := u.TokenService.Create(user.ID, name)
			if err != nil {
				data.Error = "Failed to create token"
			} else {
				http.Redirect(w, r, "/users/me", http.StatusSeeOther)
				return
			}
		}
	}

	u.Templates.Current.Execute(w, r, data)
}

func (u UsersHtml) DeleteToken(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusFound)
		return
	}
	tokenUUID := r.FormValue("uuid")
	if tokenUUID != "" && u.TokenService != nil {
		_ = u.TokenService.Delete(user.ID, tokenUUID)
	}
	http.Redirect(w, r, "/users/me", http.StatusSeeOther)
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
	UserService  *models.UserService
	TokenService *models.TokenService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := getSessionEmail(r)
		var user *models.User
		var err error

		if email != "" {
			user, err = umw.UserService.GetByEmail(email)
			if err == nil && user != nil {
				ctx := r.Context()
				ctx = context.WithUser(ctx, user)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
		}

		// Check for token in Authorization header: "Bearer <uuid>"
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			tokenStr = strings.TrimSpace(tokenStr)
			userID, err := umw.TokenService.GetUserId(tokenStr)
			if err == nil && userID > 0 {
				user, err = umw.UserService.GetByID(userID)
				if err == nil && user != nil {
					ctx := r.Context()
					ctx = context.WithUser(ctx, user)
					r = r.WithContext(ctx)
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		// No session or valid token, proceed without setting a user.
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
