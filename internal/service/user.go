package service

import (
	"context"
	"trainee-pvz/internal/models"
	"trainee-pvz/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, user models.User) error {
	return s.repo.Create(ctx, user)
}

func (s *UserService) Login(ctx context.Context, email string) (models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}
