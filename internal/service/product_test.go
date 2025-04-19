package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"trainee-pvz/internal/models"
	"trainee-pvz/internal/service"
)

type fakeMetrics struct{}

func (f *fakeMetrics) SaveEntityCount(value float64, entity string) {

}

type fakeProductRepo struct {
	addErr    error
	deleteErr error
}

func (f *fakeProductRepo) Add(ctx context.Context, p models.Product) error {
	return f.addErr
}

func (f *fakeProductRepo) DeleteLast(ctx context.Context, pvzID string) error {
	return f.deleteErr
}
func TestProductService_AddProduct_Success(t *testing.T) {
	repo := &fakeProductRepo{}
	svc := service.NewProductService(repo, &fakeMetrics{})

	err := svc.AddProduct(context.Background(), models.Product{
		ID:          "id1",
		Type:        "одежда",
		ReceptionID: "rec1",
	})
	assert.NoError(t, err)
}

func TestProductService_AddProduct_Fail(t *testing.T) {
	repo := &fakeProductRepo{addErr: errors.New("fail add")}
	svc := service.NewProductService(repo, &fakeMetrics{})

	err := svc.AddProduct(context.Background(), models.Product{
		ID:          "id2",
		Type:        "обувь",
		ReceptionID: "rec2",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fail add")
}
func TestProductService_DeleteLastProduct_Success(t *testing.T) {
	repo := &fakeProductRepo{}
	svc := service.NewProductService(repo, &fakeMetrics{})

	err := svc.DeleteLastProduct(context.Background(), "pvz1")
	assert.NoError(t, err)
}

func TestProductService_DeleteLastProduct_Fail(t *testing.T) {
	repo := &fakeProductRepo{deleteErr: errors.New("nothing to delete")}
	svc := service.NewProductService(repo, &fakeMetrics{})

	err := svc.DeleteLastProduct(context.Background(), "pvz2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nothing to delete")
}
