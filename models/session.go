package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/onehappyfellow/daebak-web/util"
)

const SessionTokenLength = 32

type Session struct {
	ID        int
	UserId    int
	Token     string // only set when creating a new session
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (s *SessionService) Create(userId int) (*Session, error) {
	token, err := util.RandomString(SessionTokenLength)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserId:    userId,
		Token:     token,
		TokenHash: s.hash(token),
	}
	row := s.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET token_hash = $2
		RETURNING id;
	`, session.UserId, session.TokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (s *SessionService) User(token string) (*User, error) {
	user := User{}
	row := s.DB.QueryRow(`
		SELECT u.id, u.email FROM users AS u
		INNER JOIN sessions AS s ON s.user_id = u.id
		WHERE s.token_hash = $1;
	`, s.hash(token))
	err := row.Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, nil
}

func (s *SessionService) Delete(token string) error {
	_, err := s.DB.Exec(`
		DELETE FROM sessions WHERE token_hash = $1;
	`, s.hash(token))
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (s *SessionService) hash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}
