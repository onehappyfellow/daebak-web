package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// Opens a SQL connection with the provided postgres config
// Callers must ensure the connection is eventually closed with db.Close()
func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

func DefaultPostresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "onehappyfellow",
		Password: "learnkorean",
		Database: "daebak",
		SSLMode:  "disable",
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (p PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", p.Host, p.Port, p.User, p.Password, p.Database, p.SSLMode)
}
