package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	UUID     uuid.UUID
	UserID   int
	Name     string
	LastUsed *time.Time
}

type TokenService struct {
	DB *sql.DB
}

func (s *TokenService) ListByUserID(userID int) ([]Token, error) {
	rows, err := s.DB.Query(`SELECT uuid, user_id, name, last_used FROM tokens WHERE user_id = $1 ORDER BY name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tokens []Token
	for rows.Next() {
		var t Token
		var uuidStr string
		if err := rows.Scan(&uuidStr, &t.UserID, &t.Name, &t.LastUsed); err != nil {
			return nil, err
		}
		t.UUID, _ = uuid.Parse(uuidStr)
		tokens = append(tokens, t)
	}
	return tokens, nil
}

func (s *TokenService) Create(userID int, name string) (*Token, error) {
	id := uuid.New()
	_, err := s.DB.Exec(
		`INSERT INTO tokens (uuid, user_id, name, last_used) VALUES ($1, $2, $3, NULL)`,
		id.String(), userID, name,
	)
	if err != nil {
		return nil, err
	}
	return &Token{
		UUID:   id,
		UserID: userID,
		Name:   name,
	}, nil
}

func (s *TokenService) GetUserId(tokenUUID string) (int, error) {
	var userID int
	err := s.DB.QueryRow(`
		UPDATE tokens
		SET last_used = NOW()
		WHERE uuid = $1
		RETURNING user_id
	`, tokenUUID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *TokenService) Delete(userID int, uuid string) error {
	_, err := s.DB.Exec(
		`DELETE FROM tokens WHERE user_id = $1 AND uuid = $2`,
		userID, uuid,
	)
	return err
}
