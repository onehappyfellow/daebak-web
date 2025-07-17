package models

import "database/sql"

type Tag struct {
	ID       int
	ParentID sql.NullInt64
	Name     string
}

type ArticleTag struct {
	ArticleID int
	TagID     int
}

type TagService struct {
	DB *sql.DB
}

func (s *TagService) GetTagByID(id int) (*Tag, error) {
	var t Tag
	err := s.DB.QueryRow(`SELECT id, parent_id, name FROM tags WHERE id = $1`, id).Scan(&t.ID, &t.ParentID, &t.Name)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Add more CRUD and ArticleTag methods as needed
