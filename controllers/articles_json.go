package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/onehappyfellow/daebak-web/models"
	"github.com/onehappyfellow/daebak-web/util"
)

type ArticlesJson struct {
	ArticleService *models.ArticleService
}

func (c ArticlesJson) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	// Set defaults if not provided
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	response, err := c.ArticleService.GetAllArticles(page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (c ArticlesJson) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := c.ArticleService.CreateArticle(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	article.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article)
}

func (c ArticlesJson) GetArticle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	article, err := c.ArticleService.GetArticle(int(id))

	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(article)
}

func (c ArticlesJson) GetArticleBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	article, err := c.ArticleService.GetArticleBySlug(slug)

	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(article)
}

func (c ArticlesJson) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	article.ID = int(id)
	article.Slug, _ = util.RandomString(models.SlugLength)

	if err := c.ArticleService.UpdateArticle(article); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(article)
}

func (c ArticlesJson) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := c.ArticleService.DeleteArticle(int(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
