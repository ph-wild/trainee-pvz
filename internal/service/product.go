package service

import (
	"context"

	"trainee-pvz/internal/models"

	//"trainee-pvz/internal/repository"
	"github.com/pkg/errors"
)

type ProductRepository interface {
	Add(ctx context.Context, product models.Product) error
	DeleteLast(ctx context.Context, pvzID string) error
}

type ProductService struct {
	repo    ProductRepository //*repository.ProductRepository
	metrics metrics
}

func NewProductService(repo ProductRepository, m metrics) *ProductService { //*repository.ProductRepository
	return &ProductService{repo: repo, metrics: m}
}

func (s *ProductService) AddProduct(ctx context.Context, p models.Product) error {
	err := s.repo.Add(ctx, p)
	if err != nil {
		return errors.Wrap(err, "can't add product")
	}

	s.metrics.SaveEntityCount(1, "product")
	return nil
}

func (s *ProductService) DeleteLastProduct(ctx context.Context, pvzID string) error {
	//ErrNoProducts, ErrNoOpenReception
	//err := s.repo.DeleteLast(ctx, receptionID)
	return s.repo.DeleteLast(ctx, pvzID)
}
