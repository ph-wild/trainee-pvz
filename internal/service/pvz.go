package service

import (
	"context"
	"time"

	er "trainee-pvz/internal/errors"
	"trainee-pvz/internal/models"
	"trainee-pvz/internal/openapi"
)

type PVZRepository interface {
	Create(ctx context.Context, pvz models.PVZ) error
	ListWithReceptionsAndProducts(
		ctx context.Context,
		start, end *time.Time,
		page, limit int,
	) ([]models.PVZWithReceptions, error)
}

type PVZService struct {
	repo PVZRepository //*repository.PVZRepository
}

func NewPVZService(repo PVZRepository) *PVZService {
	return &PVZService{repo: repo}
}

func (s *PVZService) CreatePVZ(ctx context.Context, pvz models.PVZ) error {
	switch pvz.City {
	case string(openapi.Москва), string(openapi.Казань), string(openapi.СанктПетербург):
		return s.repo.Create(ctx, pvz)
	default:
		return er.ErrUnsupportedCity
	}
}

func (s *PVZService) ListPVZ(ctx context.Context, start, end *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	return s.repo.ListWithReceptionsAndProducts(ctx, start, end, page, limit)
} //ErrNoPVZ
