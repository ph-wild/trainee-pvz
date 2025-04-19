package service

import (
	"context"

	"github.com/pkg/errors"

	er "trainee-pvz/internal/errors"
	"trainee-pvz/internal/models"
	//"trainee-pvz/internal/repository"
)

type ReceptionRepository interface {
	HasOpenReception(ctx context.Context, pvzID string) (bool, error)
	GetLastReceptionID(ctx context.Context, pvzID string) (string, error)
	Create(ctx context.Context, r models.Reception) error
	GetOpenReceptionID(ctx context.Context, pvzID string) (string, error)
	Close(ctx context.Context, id string) error
}

type ReceptionService struct {
	repo    ReceptionRepository
	metrics metrics
}

func NewReceptionService(repo ReceptionRepository, m metrics) *ReceptionService {
	return &ReceptionService{repo: repo, metrics: m}
}

func (s *ReceptionService) CreateReception(ctx context.Context, rec models.Reception) error {
	hasOpen, err := s.repo.HasOpenReception(ctx, rec.PVZID)
	if err != nil {
		return err
	}
	if hasOpen {
		return er.ErrReceptionAlreadyExists
	}

	err = s.repo.Create(ctx, rec)
	if err != nil {
		return errors.Wrap(err, "can't create reception")
	}

	s.metrics.SaveEntityCount(1, "reception")
	return nil
}

func (s *ReceptionService) CloseReception(ctx context.Context, id string) error {
	return s.repo.Close(ctx, id)
}

func (s *ReceptionService) GetLastReceptionID(ctx context.Context, pvzID string) (string, error) {
	//ErrNoOpenReception
	return s.repo.GetLastReceptionID(ctx, pvzID)
}

func (s *ReceptionService) GetOpenReceptionID(ctx context.Context, pvzID string) (string, error) {
	//ErrNoOpenReception
	return s.repo.GetOpenReceptionID(ctx, pvzID)
}
