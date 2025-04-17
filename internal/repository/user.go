package repository

import (
	"context"
	"log/slog"

	"trainee-pvz/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user models.User) error {
	query := `INSERT INTO users (id, email, password, role) VALUES (:id, :email, :password, :role)`
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		slog.Error("create user failed", slog.Any("err", err))
		return errors.Wrap(err, "repo: create user")
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	query := `SELECT id, email, password, role FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		slog.Error("get user by email failed", slog.Any("err", err))
		return user, errors.Wrap(err, "repo: get user")
	}
	return user, nil
}
