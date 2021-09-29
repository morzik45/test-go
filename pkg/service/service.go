package service

import (
	exam "github.com/morzik45/test-go"
	"github.com/morzik45/test-go/pkg/repository"
)

type Authorization interface {
	CreateUser(user exam.User) (int, error)
	GenerateToken(username, password string) (string, error)
	SessionClose(username, session_token string) (int, error)
	ParseToken(session_token string) (*exam.Authorization, error)
}

type Testing interface {
	GetAllVariants() ([]exam.Variant, error)
	GetTaskById(variantId, taskId int, username string) (exam.Task, error)
	SaveAnswer(answer exam.Answer, username string) (bool, int, error)
}

type Service struct {
	Authorization
	Testing
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Testing:       NewTasksService(repos.Testing),
	}
}
