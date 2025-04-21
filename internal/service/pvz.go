package service

import (
	"context"
	"time"

	"github.com/pkg/errors"

	er "trainee-pvz/internal/errors"
	"trainee-pvz/internal/models"
	"trainee-pvz/internal/openapi"
)

type PVZRepository interface {
	Create(ctx context.Context, pvz models.PVZ) error
	List(ctx context.Context, start, end *time.Time, page, limit int) ([]models.PVZ, error)
}

type metrics interface {
	SaveEntityCount(value float64, entity string)
}

type PVZService struct {
	repo    PVZRepository
	metrics metrics
}

func NewPVZService(repo PVZRepository, m metrics) *PVZService {
	return &PVZService{repo: repo, metrics: m}
}

func (s *PVZService) CreatePVZ(ctx context.Context, pvz models.PVZ) error {
	switch pvz.City {
	case string(openapi.Москва), string(openapi.Казань), string(openapi.СанктПетербург):
		err := s.repo.Create(ctx, pvz)
		if err != nil {
			return errors.Wrap(err, "can't create PVZ")
		}

		s.metrics.SaveEntityCount(1, "pvz")

		return nil

	default:
		return er.ErrUnsupportedCity
	}
}

func (s *PVZService) ListPVZ(ctx context.Context, start, end *time.Time, page, limit int) ([]models.PVZ, error) {
	return s.repo.List(ctx, start, end, page, limit)
}
