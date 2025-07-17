package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/onehappyfellow/daebak-web/models"
	"github.com/onehappyfellow/daebak-web/views"
)

type AdminHtml struct {
	Templates struct {
		Form views.Template
	}
	ArticleService    *models.ArticleService
	VocabularyService *models.VocabularyService
}

// Renders the form for creating a new article (no DB write)
func (c AdminHtml) NewArticleForm(w http.ResponseWriter, r *http.Request) {
	var data struct {
		models.Article
		Vocabulary interface{}
	}
	data.Vocabulary = []models.Vocabulary{}
	c.Templates.Form.Execute(w, r, data)
}

// Renders the form for editing an existing article (no DB write)
func (c AdminHtml) EditArticleForm(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}
	article, err := c.ArticleService.GetArticle(id)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}
	vocab, err := c.VocabularyService.GetVocabularyForArticle(id)
	if err != nil {
		http.Error(w, "Failed to fetch vocabulary", http.StatusInternalServerError)
		return
	}
	var data struct {
		models.Article
		Vocabulary interface{}
	}
	data.Article = *article
	data.Vocabulary = vocab
	c.Templates.Form.Execute(w, r, data)
}

// Helper to parse vocabulary IDs from form
func parseVocabularyIDs(r *http.Request) []int {
	ids := []int{}
	for _, v := range r.Form["vocabulary"] {
		if id, err := strconv.Atoi(v); err == nil {
			ids = append(ids, id)
		}
	}
	fmt.Println("Parsed vocabulary IDs:", ids)
	return ids
}
