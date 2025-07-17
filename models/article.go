package models

import (
	"database/sql"
	"math"
	"time"

	"github.com/onehappyfellow/daebak-web/util"
)

const SlugLength = 8

type Article struct {
	ID                     int        `json:"id"`
	UUID                   string     `json:"uuid"`
	Published              bool       `json:"published"`
	SourcePublished        *time.Time `json:"source_published"`
	SourceAccessed         time.Time  `json:"source_accessed"`
	SourceURL              *string    `json:"source_url"`
	SourcePublication      *string    `json:"source_publication"`
	SourceAuthor           *string    `json:"source_author"`
	Headline               string     `json:"headline"`
	HeadlineEn             *string    `json:"headline_en"`
	Content                *string    `json:"content"`
	Summary                *string    `json:"summary"`
	Context                *string    `json:"context"`
	TopikLevel             *int64     `json:"topik_level"`
	TopikLevelExplanation  *string    `json:"topik_level_explanation"`
	ComprehensionQuestions *string    `json:"comprehension_questions"`
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
			   SELECT id, uuid, published, source_published, source_accessed, source_url, source_publication, source_author, headline, headline_en, content, summary, context, topik_level, topik_level_explanation, comprehension_questions
        FROM articles WHERE id = $1;`,
		id).Scan(
		&a.ID, &a.UUID, &a.Published, &a.SourcePublished, &a.SourceAccessed, &a.SourceURL, &a.SourcePublication, &a.SourceAuthor, &a.Headline, &a.HeadlineEn, &a.Content, &a.Summary, &a.Context, &a.TopikLevel, &a.TopikLevelExplanation, &a.ComprehensionQuestions)
	return &a, err
}

// Example: get by UUID
func (s *ArticleService) GetArticleByUUID(uuid string) (*Article, error) {
	var a Article
	err := s.DB.QueryRow(`
			   SELECT id, uuid, published, source_published, source_accessed, source_url, source_publication, source_author, headline, headline_en, content, summary, context, topik_level, topik_level_explanation, comprehension_questions
			   FROM articles WHERE uuid = $1;`,
		uuid).Scan(
		&a.ID, &a.UUID, &a.Published, &a.SourcePublished, &a.SourceAccessed, &a.SourceURL, &a.SourcePublication, &a.SourceAuthor, &a.Headline, &a.HeadlineEn, &a.Content, &a.Summary, &a.Context, &a.TopikLevel, &a.TopikLevelExplanation, &a.ComprehensionQuestions)
	return &a, err
}

func (s *ArticleService) CreateArticle(a Article) (int, error) {
	// generate and set a UUID if not set
	if a.UUID == "" {
		a.UUID, _ = util.RandomString(SlugLength)
	}
	if a.SourceAccessed.IsZero() {
		a.SourceAccessed = time.Now()
	}
	var id int
	err := s.DB.QueryRow(`
			   INSERT INTO articles (uuid, published, source_published, source_accessed, source_url, source_publication, source_author, headline, headline_en, content, summary, context, topik_level, topik_level_explanation, comprehension_questions)
			   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
        RETURNING id;`,
		a.UUID, a.Published, a.SourcePublished, a.SourceAccessed, a.SourceURL, a.SourcePublication, a.SourceAuthor, a.Headline, a.HeadlineEn, a.Content, a.Summary, a.Context, a.TopikLevel, a.TopikLevelExplanation, a.ComprehensionQuestions).Scan(&id)
	return id, err
}

func (s *ArticleService) UpdateArticle(a Article) error {
	_, err := s.DB.Exec(`
        UPDATE articles 
			   SET uuid = $1, published = $2, source_published = $3, source_accessed = $4, source_url = $5, source_publication = $6, source_author = $7, headline = $8, headline_en = $9, content = $10, summary = $11, context = $12, topik_level = $13, topik_level_explanation = $14, comprehension_questions = $15
			   WHERE id = $16`,
		a.UUID, a.Published, a.SourcePublished, a.SourceAccessed, a.SourceURL, a.SourcePublication, a.SourceAuthor, a.Headline, a.HeadlineEn, a.Content, a.Summary, a.Context, a.TopikLevel, a.TopikLevelExplanation, a.ComprehensionQuestions, a.ID)
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

	offset := (page - 1) * pageSize

	// Get paginated articles
	rows, err := s.DB.Query(`
			   SELECT id, uuid, published, source_published, source_accessed, source_url, source_publication, source_author, headline, headline_en, content, summary, context, topik_level, topik_level_explanation, comprehension_questions
        FROM articles 
			   ORDER BY source_accessed DESC 
        LIMIT $1 OFFSET $2`,
		pageSize, offset)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		err := rows.Scan(
			&a.ID, &a.UUID, &a.Published, &a.SourcePublished, &a.SourceAccessed, &a.SourceURL, &a.SourcePublication, &a.SourceAuthor, &a.Headline, &a.HeadlineEn, &a.Content, &a.Summary, &a.Context, &a.TopikLevel, &a.TopikLevelExplanation, &a.ComprehensionQuestions)
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
