package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/onehappyfellow/daebak-web/models"
	"github.com/onehappyfellow/daebak-web/views"
)

type ArticlesHtml struct {
	Templates struct {
		Single views.Template
		List   views.Template
		Form   views.Template
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

func (c ArticlesHtml) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var article models.Article

	if r.Method == http.MethodPost {
		date, err := time.Parse("2006-01-02", r.FormValue("date"))
		if err != nil {
			http.Error(w, "Invalid date", http.StatusBadRequest)
			return
		}

		article.Headline = r.FormValue("headline")
		article.Content = r.FormValue("content")
		article.Date = date.UTC() // Stores the date at midnight UTC
		article.Published = r.FormValue("published") == "on"
		article.Author = r.FormValue("author")
		id, err := c.ArticleService.CreateArticle(article)
		if err != nil {
			fmt.Println("CreateArticle failed", err)
			http.Error(w, "Create article failed", http.StatusBadRequest)
			return
		}
		fmt.Printf("Created article %d\n", id)
		// TODO set success toast
		// clear form
		article = models.Article{}
	}

	if article.Date.IsZero() {
		article.Date = time.Now()
	}

	var data struct {
		models.Article
	}
	data.Article = article
	c.Templates.Form.Execute(w, r, data)
	// after this, kick off go routine that will save changes on its completion
}
