package service

import (
	"context"

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
	repo ReceptionRepository
}

func NewReceptionService(repo ReceptionRepository) *ReceptionService {
	return &ReceptionService{repo: repo}
}

func (s *ReceptionService) CreateReception(ctx context.Context, rec models.Reception) error {
	hasOpen, err := s.repo.HasOpenReception(ctx, rec.PVZID)
	if err != nil {
		return err
	}
	if hasOpen {
		return er.ErrReceptionAlreadyExists
	}
	return s.repo.Create(ctx, rec)
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
