package repository

import (
	"database/sql"

	exam "github.com/morzik45/test-go"
)

type Authorization interface {
	CreateUser(user exam.User) (int, error)
	GetUser(username, password string) (exam.User, error)
	LoginUser(username, session_token string) (int, error)
	LogoutUser(username, session_token string) (int, error)
	ParseToken(session_token string) (*exam.Authorization, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
