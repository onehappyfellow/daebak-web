package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/onehappyfellow/daebak-web/models"
	"github.com/onehappyfellow/daebak-web/views"
)

type Articles struct {
	Templates struct {
		Single views.Template
		List   views.Template
	}
	ArticleService *models.ArticleService
}

func (c Articles) Single(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> SINGLE")
	// article := r.Context().Value("article").(*Article)
	// article := r.Context().Value("article")
	// fmt.Println(article)

	if articleID := chi.URLParam(r, "articleID"); articleID != "" {
		// article, err = dbGetArticle(articleID)
		fmt.Println("articleID=", articleID)
	} else if articleSlug := chi.URLParam(r, "articleSlug"); articleSlug != "" {
		// article, err = dbGetArticleBySlug(articleSlug)
		fmt.Println("articleSlug=", articleSlug)
	}

	// var data struct {
	// 	Email string
	// }
	// data.Email = r.FormValue("email")

	// a, err := c.ArticleService.Get(slug)
	// if err != nil {
	// 	if err == ErrNotFound {
	// 		// redirect to 404
	// 		return
	// 	}
	// 	// 500 error
	// 	return
	// }
	c.Templates.Single.Execute(w, r, nil)
}

func (c Articles) Trending(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> TRENDING")
	// var data struct {
	// 	Email string
	// }
	// data.Email = r.FormValue("email")
	c.Templates.List.Execute(w, r, nil)
}

type ArticleMiddleware struct {
	ArticleService *models.ArticleService
}

func (mw ArticleMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cookie, err := r.Cookie(CookieSession)
		// if err != nil {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }
		// user, err := mw.SessionService.User(cookie.Value)
		// if err != nil {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }
		// ctx := context.WithUser(r.Context(), user)
		// next.ServeHTTP(w, r.WithContext(ctx))

		if articleID := chi.URLParam(r, "articleID"); articleID != "" {
			// article, err = dbGetArticle(articleID)
			// article, err = c.ArticleService.Get()
			fmt.Println("articleID=", articleID)
		} else if articleSlug := chi.URLParam(r, "articleSlug"); articleSlug != "" {
			// article, err = dbGetArticleBySlug(articleSlug)
			fmt.Println("articleSlug=", articleSlug)
		} else {
			fmt.Println("nothing")
			// render.Render(w, r, ErrNotFound)
			// http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		ctx := context.WithValue(r.Context(), "article", "jonathan")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// var article *Article
		// var err error

		if articleID := chi.URLParam(r, "articleID"); articleID != "" {
			// article, err = dbGetArticle(articleID)
			// article, err = c.ArticleService.Get()
			fmt.Println("articleID=", articleID)
		} else if articleSlug := chi.URLParam(r, "articleSlug"); articleSlug != "" {
			// article, err = dbGetArticleBySlug(articleSlug)
			fmt.Println("articleSlug=", articleSlug)
		} else {
			fmt.Println("nothing")
			// render.Render(w, r, ErrNotFound)
			return
		}
		// if err != nil {
		// 	// render.Render(w, r, ErrNotFound)
		// 	return
		// }

		// ctx := context.WithValue(r.Context(), "article", article)
		ctx := context.WithValue(r.Context(), "article", "jonathan")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
