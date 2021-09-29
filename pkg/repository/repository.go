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

type Testing interface {
	StartTest(userId, variantId int) (int, error)
	GetAllVariants() ([]exam.Variant, error)
	GetTaskById(variantId, taskId int, username string) (exam.Task, error)
	SaveAnswer(answer exam.Answer, username string) (bool, int, error)
}

type Repository struct {
	Authorization
	Testing
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Testing:       NewTaskPostgres(db),
	}
}
