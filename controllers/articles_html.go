package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/onehappyfellow/daebak-web/models"
	"github.com/onehappyfellow/daebak-web/views"
)

type ArticlesHtml struct {
	Templates struct {
		Single views.Template
		List   views.Template
	}
	ArticleService *models.ArticleService
}

func (c ArticlesHtml) Single(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	article, err := c.ArticleService.GetArticleBySlug(slug)

	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}
	c.Templates.Single.Execute(w, r, article)
}

func (c ArticlesHtml) Trending(w http.ResponseWriter, r *http.Request) {

	page, err := c.ArticleService.GetAllArticles(1, 10)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	var data struct {
		Title    string
		Articles []models.Article
	}
	data.Title = "Trndin"
	data.Articles = page.Articles
	c.Templates.List.Execute(w, r, data)
}
