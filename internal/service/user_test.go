package service_test

import (
	"context"
	"errors"
	"testing"

	"trainee-pvz/internal/models"
	"trainee-pvz/internal/service"

	"github.com/stretchr/testify/assert"
)

type fakeUserRepo struct {
	createErr error
	getUser   models.User
	getErr    error
}

func (f *fakeUserRepo) Create(ctx context.Context, user models.User) error {
	return f.createErr
}

func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	if f.getErr != nil {
		return models.User{}, f.getErr
	}
	return f.getUser, nil
}

func TestUserService_Register_Success(t *testing.T) {
	repo := &fakeUserRepo{}
	svc := service.NewUserService(repo)

	err := svc.Register(context.Background(), models.User{
		ID:    "id1",
		Email: "test@example.com",
	})
	assert.NoError(t, err)
}

func TestUserService_Register_Fail(t *testing.T) {
	repo := &fakeUserRepo{createErr: errors.New("duplicate")}
	svc := service.NewUserService(repo)

	err := svc.Register(context.Background(), models.User{
		ID:    "id2",
		Email: "test@example.com",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate")
}

func TestUserService_Login_Success(t *testing.T) {
	repo := &fakeUserRepo{
		getUser: models.User{
			ID:    "id3",
			Email: "user@example.com",
			Role:  "employee",
		},
	}
	svc := service.NewUserService(repo)

	user, err := svc.Login(context.Background(), "user@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "user@example.com", user.Email)
}

func TestUserService_Login_NotFound(t *testing.T) {
	repo := &fakeUserRepo{getErr: errors.New("not found")}
	svc := service.NewUserService(repo)

	_, err := svc.Login(context.Background(), "notfound@example.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
