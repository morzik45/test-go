package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	usersTable          = "users"
	authorizationsTable = "authorizations"
	variantsTable       = "variants"
	testsTable          = "tests"
	tasksTable          = "tasks"
	userAnswersTable    = "user_answers"
	resultsTable        = "results"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
