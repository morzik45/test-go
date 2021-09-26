package service

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os/exec"
	"strings"

	exam "github.com/morzik45/test-go"
	"github.com/morzik45/test-go/pkg/repository"
)

const (
	salt = "hjqrhjqw124617ajfhajs"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user exam.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) SessionClose(username, session_token string) (int, error) {
	return s.repo.LogoutUser(username, session_token)
}

func (s *AuthService) ParseToken(session_token string) (*exam.Authorization, error) {
	return s.repo.ParseToken(session_token)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	uuid, err := exec.Command("uuidgen").Output()
	newSessionToken := strings.TrimRight(string(uuid), "\n")
	if err != nil {
		return "", err
	}

	_, err = s.repo.LoginUser(username, newSessionToken)
	if err != nil {
		return "", err
	}

	log.Printf("User %s authorised.", user.Username)
	return newSessionToken, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
