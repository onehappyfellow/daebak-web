package models

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int
	Email             string
	PasswordHash      string
	ResetToken        sql.NullString
	ResetTokenExpires sql.NullTime
	CreatedAt         time.Time
}

type UserService struct {
	DB *sql.DB
}

func (s *UserService) CreateUser(email, password string) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	var id int
	err = s.DB.QueryRow(`
		INSERT INTO users (email, password_hash, created_at)
		VALUES ($1, $2, NOW()) RETURNING id;`,
		email, string(hash)).Scan(&id)
	return id, err
}

func (s *UserService) Authenticate(email, password string) (*User, error) {
	var u User
	err := s.DB.QueryRow(`
		SELECT id, email, password_hash, created_at FROM users WHERE email = $1;`,
		email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return nil, sql.ErrNoRows
	}
	return &u, nil
}

func (s *UserService) SetResetToken(email string) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	hash := sha256.Sum256([]byte(token))
	hashHex := hex.EncodeToString(hash[:])
	expires := time.Now().UTC().Add(1 * time.Hour)
	_, err = s.DB.Exec(`
		UPDATE users SET reset_token = $1, reset_token_expires = $2 WHERE email = $3;`,
		hashHex, expires, email)
	return token, err
}

func (s *UserService) GetByResetToken(token string) (*User, error) {
	var u User
	hash := sha256.Sum256([]byte(token))
	hashHex := hex.EncodeToString(hash[:])
	err := s.DB.QueryRow(`
		SELECT id, email, password_hash, reset_token_expires FROM users WHERE reset_token = $1;`,
		hashHex).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.ResetTokenExpires)
	if err != nil {
		return nil, err
	}
	// Always compare in UTC
	now := time.Now().UTC()
	if u.ResetTokenExpires.Valid && u.ResetTokenExpires.Time.After(now) {
		return &u, nil
	}
	return nil, sql.ErrNoRows
}

func (s *UserService) ResetPassword(id int, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(`
		UPDATE users SET password_hash = $1, reset_token = NULL, reset_token_expires = NULL WHERE id = $2;`,
		string(hash), id)
	return err
}

func (s *UserService) GetByEmail(email string) (*User, error) {
	var u User
	err := s.DB.QueryRow(`
		SELECT id, email, password_hash, created_at FROM users WHERE email = $1;`,
		email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserService) GetByID(id int) (*User, error) {
	var u User
	err := s.DB.QueryRow(`
		SELECT id, email, password_hash, created_at FROM users WHERE id = $1;`,
		id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
