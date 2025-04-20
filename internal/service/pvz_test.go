package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"trainee-pvz/internal/models"
	"trainee-pvz/internal/openapi"
	"trainee-pvz/internal/service"
)

type fakePVZRepo struct {
	createErr error
	listErr   error
	data      []models.PVZWithReceptions
}

func (f *fakePVZRepo) Create(ctx context.Context, pvz models.PVZ) error {
	return f.createErr
}

func (f *fakePVZRepo) ListWithReceptionsAndProducts(
	ctx context.Context,
	start, end *time.Time,
	page, limit int,
) ([]models.PVZWithReceptions, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.data, nil
}

func TestPVZService_CreatePVZ_Success(t *testing.T) {
	repo := &fakePVZRepo{}
	svc := service.NewPVZService(repo, &fakeMetrics{})

	err := svc.CreatePVZ(context.Background(), models.PVZ{
		ID:   "1",
		City: string(openapi.Москва),
	})
	assert.NoError(t, err)
}

func TestPVZService_CreatePVZ_UnsupportedCity(t *testing.T) {
	repo := &fakePVZRepo{}
	svc := service.NewPVZService(repo, &fakeMetrics{})

	err := svc.CreatePVZ(context.Background(), models.PVZ{
		ID:   "2",
		City: "Новосибирск",
	})
	assert.Error(t, err)
	assert.Equal(t, "unsupported city", err.Error())
}

func TestPVZService_CreatePVZ_RepoError(t *testing.T) {
	repo := &fakePVZRepo{createErr: errors.New("db error")}
	svc := service.NewPVZService(repo, &fakeMetrics{})

	err := svc.CreatePVZ(context.Background(), models.PVZ{
		ID:   "3",
		City: string(openapi.Казань),
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestPVZService_ListPVZ_Success(t *testing.T) {
	repo := &fakePVZRepo{
		data: []models.PVZWithReceptions{
			{
				PVZ: models.PVZ{ID: "1", City: "Казань"},
			},
		},
	}
	svc := service.NewPVZService(repo, &fakeMetrics{})

	result, err := svc.ListPVZ(context.Background(), nil, nil, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Казань", result[0].PVZ.City)
}

func TestPVZService_ListPVZ_Empty(t *testing.T) {
	repo := &fakePVZRepo{data: []models.PVZWithReceptions{}}
	svc := service.NewPVZService(repo, &fakeMetrics{})

	result, err := svc.ListPVZ(context.Background(), nil, nil, 1, 10)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestPVZService_ListPVZ_Error(t *testing.T) {
	repo := &fakePVZRepo{listErr: errors.New("list fail")}
	svc := service.NewPVZService(repo, &fakeMetrics{})

	result, err := svc.ListPVZ(context.Background(), nil, nil, 1, 10)
	assert.Error(t, err)
	assert.Nil(t, result)
}
