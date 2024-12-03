package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/onehappyfellow/daebak-web/util"
)

const (
	ResetDuration = 10 * time.Minute
	BytesPerToken = 8
)

var (
	ErrTokenInvalid = errors.New("password_reset: token is invalid")
	ErrTokenExpired = errors.New("password_reset: token is expired")
)

type PasswordReset struct {
	ID        int
	UserId    int
	Token     string // only set when creating a new instance
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
}

func (s *PasswordResetService) Create(email string) (*PasswordReset, error) {
	// get userId from email
	email = strings.ToLower(email)
	var userId int
	row := s.DB.QueryRow(`
		SELECT id FROM users WHERE email = $1;
	`, email)
	err := row.Scan(&userId)
	if err != nil {
		// TODO consider returning a specific error when user does not exist
		return nil, fmt.Errorf("create: %w", err)
	}
	// create token for user
	token, err := util.RandomString(BytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	tokenHash := s.hash(token)
	reset := PasswordReset{
		UserId:    userId,
		Token:     token,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(ResetDuration),
	}
	// save reset in db
	row = s.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id;
	`, reset.UserId, reset.TokenHash, reset.ExpiresAt)
	err = row.Scan(&reset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &reset, nil
}

func (s *PasswordResetService) Consume(token string) (*User, error) {
	var user User
	var expiresAt time.Time
	row := s.DB.QueryRow(`
		SELECT users.id, email, expires_at 
		FROM users
		JOIN password_resets ON password_resets.user_id = users.id
		WHERE password_resets.token_hash = $1;
	`, s.hash(token))
	err := row.Scan(&user.ID, &user.Email, &expiresAt)
	if err != nil {
		return nil, ErrTokenInvalid
	}
	if time.Now().After(expiresAt) {
		return nil, ErrTokenExpired
	}
	fmt.Println("CONSUME CONTINUED")
	// _, err = s.DB.Exec(`
	// 	DELETE FROM password_resets WHERE user_id = $1;
	// `, user.ID)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	return &user, nil
}

func (s *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
