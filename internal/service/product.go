package service

import (
	"context"

	"trainee-pvz/internal/models"
	"trainee-pvz/internal/repository"
	//er "github.com/pkg/errors"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) AddProduct(ctx context.Context, p models.Product) error {
	return s.repo.Add(ctx, p)
}

func (s *ProductService) DeleteLastProduct(ctx context.Context, pvzID string) error {
	//ErrNoProducts, ErrNoOpenReception
	//err := s.repo.DeleteLast(ctx, receptionID)
	return s.repo.DeleteLast(ctx, pvzID)
}
