package models

import (
	"database/sql"
	"errors"
)

const SlugLength = 8

var (
	ErrNotFound = errors.New("article: not found")
)

type Article struct {
	ID       int
	Slug     string
	Headline string
	Body     string
}

type ArticleService struct {
	DB *sql.DB
}

// func (s *ArticleService) Create(email, password string) (*Article, error) {
// 	a := Article{
// 	}
// 	row := s.DB.QueryRow(`
// 		INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id;
// 	`, user.Email, passwordHash)
// 	err = row.Scan(&user.ID)
// 	if err != nil {
// 		var pgError *pgconn.PgError
// 		if errors.As(err, &pgError) {
// 			if pgError.Code == pgerrcode.UniqueViolation {
// 				return nil, ErrEmailTaken
// 			}
// 		}
// 		return nil, fmt.Errorf("create user: %w", err)
// 	}
// 	return &user, nil
// }

func (s *ArticleService) Get(id int) (*Article, error) {
	a := Article{
		Slug:     "s-a-ge-help",
		Headline: "Article Headline",
		Body:     `This is the article text.\nIt could be very long.`,
	}
	// row := s.DB.QueryRow(`
	// 	SELECT id, password_hash FROM articles WHERE slug = $1;
	// `, a.Slug)
	// err := row.Scan(&user.ID, &user.PasswordHash)
	return &a, nil

}

func (s *ArticleService) GetBySlug(slug string) (*Article, error) {
	a := Article{
		Slug:     slug,
		Headline: "Article Headline",
		Body:     `This is the article text.\nIt could be very long.`,
	}
	// row := s.DB.QueryRow(`
	// 	SELECT id, password_hash FROM articles WHERE slug = $1;
	// `, a.Slug)
	// err := row.Scan(&user.ID, &user.PasswordHash)
	return &a, nil

}
