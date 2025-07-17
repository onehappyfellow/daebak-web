package models

import (
	"database/sql"
	"math"
)

type Vocabulary struct {
	ID          int     `json:"id"`
	Word        string  `json:"word"`
	Definition  *string `json:"definition"`
	Examples    *string `json:"examples"`
	Translation *string `json:"translation_en"`
}

type VocabularyPaginatedResponse struct {
	Vocabulary  []Vocabulary `json:"vocabulary"`
	TotalCount  int          `json:"total_count"`
	CurrentPage int          `json:"current_page"`
	TotalPages  int          `json:"total_pages"`
	PageSize    int          `json:"page_size"`
}

type VocabularyService struct {
	DB *sql.DB
}

func (s *VocabularyService) GetVocabularyByID(id int) (*Vocabulary, error) {
	var v Vocabulary
	err := s.DB.QueryRow(`SELECT id, word, definition, examples, translation_en FROM vocabulary WHERE id = $1`, id).Scan(
		&v.ID, &v.Word, &v.Definition, &v.Examples, &v.Translation)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *VocabularyService) GetOrCreateVocabulary(word string) (*Vocabulary, error) {
	var v Vocabulary
	err := s.DB.QueryRow(`SELECT id, word, definition, examples, translation_en FROM vocabulary WHERE word = $1`, word).Scan(
		&v.ID, &v.Word, &v.Definition, &v.Examples, &v.Translation)
	if err == nil {
		return &v, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}
	// Not found, create it
	v.Word = word
	deff := "incomplete: todo call tool"
	v.Definition = &deff
	createErr := s.DB.QueryRow(`INSERT INTO vocabulary (word, definition) VALUES ($1, $2) RETURNING id`, v.Word, v.Definition).Scan(&v.ID)
	if createErr != nil {
		return nil, createErr
	}
	return &v, nil
}

func (s *VocabularyService) ListVocabulary(page, pageSize int) (VocabularyPaginatedResponse, error) {
	var response VocabularyPaginatedResponse
	var totalCount int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM vocabulary").Scan(&totalCount)
	if err != nil {
		return response, err
	}
	offset := (page - 1) * pageSize
	rows, err := s.DB.Query(`SELECT id, word, definition, examples, translation_en FROM vocabulary ORDER BY id DESC LIMIT $1 OFFSET $2`, pageSize, offset)
	if err != nil {
		return response, err
	}
	defer rows.Close()
	var vocabList []Vocabulary
	for rows.Next() {
		var v Vocabulary
		err := rows.Scan(&v.ID, &v.Word, &v.Definition, &v.Examples, &v.Translation)
		if err != nil {
			return response, err
		}
		vocabList = append(vocabList, v)
	}
	response.Vocabulary = vocabList
	response.TotalCount = totalCount
	response.CurrentPage = page
	response.PageSize = pageSize
	response.TotalPages = int(math.Ceil(float64(totalCount) / float64(pageSize)))
	return response, nil
}

func (s *VocabularyService) GetVocabularyForArticle(articleID int) ([]Vocabulary, error) {
	rows, err := s.DB.Query(`
        SELECT v.id, v.word, v.definition, v.examples, v.translation_en
        FROM vocabulary v
        INNER JOIN article_vocabulary av ON av.vocabulary_id = v.id
        WHERE av.article_id = $1
        ORDER BY v.word ASC
    `, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vocabList []Vocabulary
	for rows.Next() {
		var v Vocabulary
		err := rows.Scan(&v.ID, &v.Word, &v.Definition, &v.Examples, &v.Translation)
		if err != nil {
			return nil, err
		}
		vocabList = append(vocabList, v)
	}
	return vocabList, nil
}

func (s *VocabularyService) CreateVocabulary(v Vocabulary) (int, error) {
	var id int
	err := s.DB.QueryRow(`INSERT INTO vocabulary (word, definition, examples, translation_en) VALUES ($1, $2, $3, $4) RETURNING id`, v.Word, v.Definition, v.Examples, v.Translation).Scan(&id)
	return id, err
}

func (s *VocabularyService) UpdateVocabulary(v Vocabulary) error {
	_, err := s.DB.Exec(`UPDATE vocabulary SET word = $1, definition = $2, examples = $3, translation_en = $4 WHERE id = $5`, v.Word, v.Definition, v.Examples, v.Translation, v.ID)
	return err
}

func (s *VocabularyService) DeleteVocabulary(id int) error {
	_, err := s.DB.Exec(`DELETE FROM vocabulary WHERE id = $1`, id)
	return err
}

func (s *VocabularyService) SetArticleVocabulary(articleID int, vocabIDs []int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove all existing associations
	_, err = tx.Exec(`DELETE FROM article_vocabulary WHERE article_id = $1`, articleID)
	if err != nil {
		return err
	}

	// Insert new associations, ignoring duplicates
	stmt, err := tx.Prepare(`INSERT INTO article_vocabulary (article_id, vocabulary_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, vocabID := range vocabIDs {
		if _, err := stmt.Exec(articleID, vocabID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
