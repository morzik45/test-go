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

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
