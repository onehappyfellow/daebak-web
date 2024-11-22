package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (s *UserService) Create(email, password string) (*User, error) {
	// normalize email
	email = strings.ToLower(email)
	email = strings.TrimSpace(email)
	// hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user password hash: %w", err)
	}
	passwordHash := string(hashedBytes)
	// create user
	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}
	row := s.DB.QueryRow(`
		INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id;
	`, email, passwordHash)
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func (s *UserService) Authenticate(email, password string) (*User, error) {
	// normalize email
	email = strings.ToLower(email)
	email = strings.TrimSpace(email)
	// query user from db
	user := User{
		Email: email,
	}
	row := s.DB.QueryRow(`
		SELECT id, password_hash FROM users WHERE email = $1;
	`, email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	// compare password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	return &user, nil
}
