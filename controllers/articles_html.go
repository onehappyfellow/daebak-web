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
	uuid := chi.URLParam(r, "slug")
	article, err := c.ArticleService.GetArticleByUUID(uuid)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}
	var data struct {
		Article models.Article
	}
	data.Article = *article
	c.Templates.Single.Execute(w, r, data)
}

func (c ArticlesHtml) Home(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title    string
		Articles []models.Article
	}
	page, err := c.ArticleService.GetAllArticles(1, 10)
	if err != nil {
		data.Articles = []models.Article{}
	} else {
		data.Articles = page.Articles
	}

	data.Title = "Home"
	c.Templates.List.Execute(w, r, data)
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
