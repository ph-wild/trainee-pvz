package service_test

import (
	"context"
	"errors"
	"testing"

	"trainee-pvz/internal/models"
	"trainee-pvz/internal/service"

	"github.com/stretchr/testify/assert"
)

type fakeReceptionRepo struct {
	hasOpen          bool
	hasOpenErr       error
	createErr        error
	closeErr         error
	lastReceptionID  string
	lastReceptionErr error
	openReceptionID  string
	openReceptionErr error
}

func (f *fakeReceptionRepo) HasOpenReception(ctx context.Context, pvzID string) (bool, error) {
	return f.hasOpen, f.hasOpenErr
}

func (f *fakeReceptionRepo) Create(ctx context.Context, r models.Reception) error {
	return f.createErr
}

func (f *fakeReceptionRepo) Close(ctx context.Context, id string) error {
	return f.closeErr
}

func (f *fakeReceptionRepo) GetLastReceptionID(ctx context.Context, pvzID string) (string, error) {
	return f.lastReceptionID, f.lastReceptionErr
}

func (f *fakeReceptionRepo) GetOpenReceptionID(ctx context.Context, pvzID string) (string, error) {
	return f.openReceptionID, f.openReceptionErr
}

func TestReceptionService_CreateReception_Success(t *testing.T) {
	repo := &fakeReceptionRepo{hasOpen: false}
	svc := service.NewReceptionService(repo, &fakeMetrics{})

	err := svc.CreateReception(context.Background(), models.Reception{
		ID:    "r1",
		PVZID: "pvz1",
	})
	assert.NoError(t, err)
}

func TestReceptionService_CreateReception_AlreadyExists(t *testing.T) {
	repo := &fakeReceptionRepo{hasOpen: true}
	svc := service.NewReceptionService(repo, &fakeMetrics{})

	err := svc.CreateReception(context.Background(), models.Reception{
		ID:    "r2",
		PVZID: "pvz2",
	})
	assert.Error(t, err)
	assert.Equal(t, "reception already exists", err.Error())

}

func TestReceptionService_CreateReception_HasOpenErr(t *testing.T) {
	repo := &fakeReceptionRepo{hasOpenErr: errors.New("db error")}
	svc := service.NewReceptionService(repo, &fakeMetrics{})

	err := svc.CreateReception(context.Background(), models.Reception{
		ID:    "r3",
		PVZID: "pvz3",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestReceptionService_CloseReception_Success(t *testing.T) {
	repo := &fakeReceptionRepo{}
	svc := service.NewReceptionService(repo, &fakeMetrics{})

	err := svc.CloseReception(context.Background(), "reception-id")
	assert.NoError(t, err)
}

func TestReceptionService_GetLastReceptionID_Success(t *testing.T) {
	repo := &fakeReceptionRepo{lastReceptionID: "last-id"}
	svc := service.NewReceptionService(repo, &fakeMetrics{})

	id, err := svc.GetLastReceptionID(context.Background(), "pvz-id")
	assert.NoError(t, err)
	assert.Equal(t, "last-id", id)
}

func TestReceptionService_GetOpenReceptionID_Success(t *testing.T) {
	repo := &fakeReceptionRepo{openReceptionID: "open-id"}
	svc := service.NewReceptionService(repo, &fakeMetrics{})

	id, err := svc.GetOpenReceptionID(context.Background(), "pvz-id")
	assert.NoError(t, err)
	assert.Equal(t, "open-id", id)
}

func TestReceptionService_GetOpenReceptionID_Error(t *testing.T) {
	repo := &fakeReceptionRepo{openReceptionErr: errors.New("no open")}
	svc := service.NewReceptionService(repo, &fakeMetrics{})

	_, err := svc.GetOpenReceptionID(context.Background(), "pvz-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no open")
}
