package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/onehappyfellow/daebak-web/models"
)

type VocabularyJson struct {
	VocabularyService *models.VocabularyService
}

func (c VocabularyJson) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	response, err := c.VocabularyService.ListVocabulary(page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(response)
}

func (c VocabularyJson) Create(w http.ResponseWriter, r *http.Request) {
	var vocab models.Vocabulary
	if err := json.NewDecoder(r.Body).Decode(&vocab); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := c.VocabularyService.CreateVocabulary(vocab)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vocab.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vocab)
}

func (c VocabularyJson) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var vocab models.Vocabulary
	if err := json.NewDecoder(r.Body).Decode(&vocab); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	vocab.ID = id
	if err := c.VocabularyService.UpdateVocabulary(vocab); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(vocab)
}

func (c VocabularyJson) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if err := c.VocabularyService.DeleteVocabulary(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
