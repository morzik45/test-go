package service

import (
	exam "github.com/morzik45/test-go"
	"github.com/morzik45/test-go/pkg/repository"
)

type TasksService struct {
	repo repository.Testing
}

func NewTasksService(repo repository.Testing) *TasksService {
	return &TasksService{repo: repo}
}

func (s *TasksService) GetAllVariants() ([]exam.Variant, error) {
	return s.repo.GetAllVariants()
}
