package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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

func (c AdminHtml) CreateArticle(w http.ResponseWriter, r *http.Request) {
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
		// Parse vocabulary IDs from form
		vocabIDs := parseVocabularyIDs(r)
		if len(vocabIDs) > 0 {
			if err := c.VocabularyService.SetArticleVocabulary(id, vocabIDs); err != nil {
				fmt.Println("SetArticleVocabulary failed", err)
				http.Error(w, "Failed to set vocabulary", http.StatusInternalServerError)
				return
			}
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
		Vocabulary interface{}
	}
	data.Vocabulary = []models.Vocabulary{} // Initialize vocabulary as empty slice
	data.Article = article
	c.Templates.Form.Execute(w, r, data)
	// after this, kick off go routine that will save changes on its completion
}

func (c AdminHtml) EditArticle(w http.ResponseWriter, r *http.Request) {
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

	// Fetch vocabulary for the article
	vocab, err := c.VocabularyService.GetVocabularyForArticle(id)
	if err != nil {
		http.Error(w, "Failed to fetch vocabulary", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		date, err := time.Parse("2006-01-02", r.FormValue("date"))
		if err != nil {
			http.Error(w, "Invalid date", http.StatusBadRequest)
			return
		}
		article.Headline = r.FormValue("headline")
		article.Content = r.FormValue("content")
		article.Date = date.UTC()
		article.Published = r.FormValue("published") == "on"
		article.Author = r.FormValue("author")
		err = c.ArticleService.UpdateArticle(*article)
		if err != nil {
			fmt.Println("UpdateArticle failed", err)
			http.Error(w, "Update article failed", http.StatusBadRequest)
			return
		}
		// Parse vocabulary IDs from form
		vocabIDs := parseVocabularyIDs(r)
		if err := c.VocabularyService.SetArticleVocabulary(article.ID, vocabIDs); err != nil {
			fmt.Println("SetArticleVocabulary failed", err)
			http.Error(w, "Failed to set vocabulary", http.StatusInternalServerError)
			return
		}
		vocab, err = c.VocabularyService.GetVocabularyForArticle(id)
		if err != nil {
			http.Error(w, "Failed to fetch vocabulary", http.StatusInternalServerError)
			return
		}
		// Optionally redirect or set a success message
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
