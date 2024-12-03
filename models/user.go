package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"golang.org/x/crypto/bcrypt"
)

const MinPasswordLength = 8

var (
	ErrEmailTaken       = errors.New("user: email address already in use")
	ErrPasswordInsecure = errors.New("user: insecure password not allowed")
	ErrInvalidAuth      = errors.New("user: invalid authentication")
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
	if !s.isPasswordSecure(password) {
		return nil, ErrPasswordInsecure
	}
	// hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user password hash: %w", err)
	}
	passwordHash := string(hashedBytes)
	// create user
	user := User{
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: passwordHash,
	}
	row := s.DB.QueryRow(`
		INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id;
	`, user.Email, passwordHash)
	err = row.Scan(&user.ID)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return nil, ErrEmailTaken
			}
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func (s *UserService) Authenticate(email, password string) (*User, error) {
	// query user from db
	user := User{
		Email: strings.ToLower(strings.TrimSpace(email)),
	}
	row := s.DB.QueryRow(`
		SELECT id, password_hash FROM users WHERE email = $1;
	`, user.Email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		// TODO check if errors.Is(err, pgx.ErrNoRows) {}
		// user is not found
		return nil, ErrInvalidAuth
	}

	// compare password hashes
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// password does not match
		return nil, ErrInvalidAuth
	}
	return &user, nil
}

func (s *UserService) UpdatePassword(user *User, password string) error {
	if !s.isPasswordSecure(password) {
		return ErrPasswordInsecure
	}
	// hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	passwordHash := string(hashedBytes)
	_, err = s.DB.Exec(`
		UPDATE users SET password_hash = $1 WHERE id = $2;
	`, passwordHash, user.ID)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

func (s *UserService) isPasswordSecure(password string) bool {
	return len(password) >= MinPasswordLength
}
