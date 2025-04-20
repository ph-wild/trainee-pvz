package service

import (
	"context"
	"trainee-pvz/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	GetByEmail(ctx context.Context, email string) (models.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, user models.User) error {
	return s.repo.Create(ctx, user)
}

func (s *UserService) Login(ctx context.Context, email string) (models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}
