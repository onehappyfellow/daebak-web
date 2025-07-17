package models

import "database/sql"

type Grammar struct {
	ID                int
	Title             string
	Explanation       sql.NullString
	ExplanationShort  sql.NullString
	Examples          sql.NullString
	Practice          sql.NullString
}

type ArticleGrammar struct {
	GrammarID      int
	ArticleID      int
	ArticleExample sql.NullString
}

type GrammarService struct {
	DB *sql.DB
}

// Add CRUD methods as needed, e.g.:
func (s *GrammarService) GetGrammarByID(id int) (*Grammar, error) {
	var g Grammar
	err := s.DB.QueryRow(`SELECT id, title, explanation, explanation_short, examples, practice FROM grammar WHERE id = $1`, id).Scan(
		&g.ID, &g.Title, &g.Explanation, &g.ExplanationShort, &g.Examples, &g.Practice)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// Add more methods for ArticleGrammar as needed
