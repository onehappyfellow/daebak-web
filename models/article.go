package models

import (
	"database/sql"
	"math"
	"time"

	"github.com/onehappyfellow/daebak-web/util"
)

const SlugLength = 8

type Article struct {
	ID        int       `json:"id"`
	Slug      string    `json:"slug"`
	Headline  string    `json:"headline"`
	Content   string    `json:"content"`
	Date      time.Time `json:"date"`
	Published bool      `json:"published"`
	Author    string    `json:"author"`
}

type PaginatedResponse struct {
	Articles    []Article `json:"articles"`
	TotalCount  int       `json:"total_count"`
	CurrentPage int       `json:"current_page"`
	TotalPages  int       `json:"total_pages"`
	PageSize    int       `json:"page_size"`
}

type ArticleService struct {
	DB *sql.DB
}

func (s *ArticleService) GetArticle(id int) (*Article, error) {
	var a Article
	err := s.DB.QueryRow(`
        SELECT id, headline, slug, content, date, published, author 
        FROM articles WHERE id = $1;`,
		id).Scan(&a.ID, &a.Headline, &a.Slug, &a.Content, &a.Date, &a.Published, &a.Author)
	return &a, err
}

func (s *ArticleService) GetArticleBySlug(slug string) (*Article, error) {
	var a Article
	err := s.DB.QueryRow(`
        SELECT id, headline, slug, content, date, published, author 
        FROM articles WHERE slug = $1;`,
		slug).Scan(&a.ID, &a.Headline, &a.Slug, &a.Content, &a.Date, &a.Published, &a.Author)
	return &a, err
}

func (s *ArticleService) CreateArticle(a Article) (int, error) {
	// generate and set a slug
	a.Slug, _ = util.RandomString(SlugLength)

	// if date is not set, use today
	if a.Date.IsZero() {
		a.Date = time.Now()
	}

	var id int
	err := s.DB.QueryRow(`
        INSERT INTO articles (headline, slug, content, date, published, author)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id;`,
		a.Headline, a.Slug, a.Content, a.Date, a.Published, a.Author).Scan(&id)
	return id, err
}

func (s *ArticleService) UpdateArticle(a Article) error {
	_, err := s.DB.Exec(`
        UPDATE articles 
        SET headline = $1, slug = $2, content = $3, published = $4, author = $5
        WHERE id = $6`,
		a.Headline, a.Slug, a.Content, a.Published, a.Author, a.ID)
	return err
}

func (s *ArticleService) DeleteArticle(id int) error {
	_, err := s.DB.Exec(`DELETE FROM articles WHERE id = $1`, id)
	return err
}

func (s *ArticleService) GetAllArticles(page, pageSize int) (PaginatedResponse, error) {
	var response PaginatedResponse

	// TODO don't count everything every time
	var totalCount int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM articles").Scan(&totalCount)
	if err != nil {
		return response, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get paginated articles
	rows, err := s.DB.Query(`
        SELECT id, headline, slug, content, date, published, author 
        FROM articles 
        ORDER BY date DESC 
        LIMIT $1 OFFSET $2`,
		pageSize, offset)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		err := rows.Scan(&a.ID, &a.Headline, &a.Slug, &a.Content, &a.Date, &a.Published, &a.Author)
		if err != nil {
			return response, err
		}
		articles = append(articles, a)
	}

	response.Articles = articles
	response.TotalCount = totalCount
	response.CurrentPage = page
	response.PageSize = pageSize
	response.TotalPages = int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return response, nil
}
