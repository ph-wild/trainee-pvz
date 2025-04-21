package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"trainee-pvz/internal/models"
)

type PVZRepository struct {
	db *sqlx.DB
}

func NewPVZRepository(db *sqlx.DB) *PVZRepository {
	return &PVZRepository{db: db}
}

func (r *PVZRepository) Create(ctx context.Context, pvz models.PVZ) error {
	query := `INSERT INTO pvz (id, city, registration_date) VALUES (:id, :city, :registration_date)`
	_, err := r.db.NamedExecContext(ctx, query, pvz)
	if err != nil {
		slog.Error("create pvz failed", slog.Any("err", err))
		return errors.Wrap(err, "repo: create pvz")
	}
	return nil
}

func (r *PVZRepository) List(ctx context.Context, start, end *time.Time, page, limit int) ([]models.PVZ, error) {
	offset := (page - 1) * limit
	var pvzList []models.PVZ

	slog.Info("start", slog.Any("time", start))
	slog.Info("end", slog.Any("time", end))

	pvzQuery := `
		SELECT id, city, registration_date 
		FROM pvz
		WHERE ($1::timestamptz IS NULL OR registration_date >= $1) AND ($2::timestamptz IS NULL OR registration_date <= $2)
		ORDER BY registration_date DESC
		OFFSET $3 LIMIT $4;
	`
	err := r.db.SelectContext(ctx, &pvzList, pvzQuery, start, end, offset, limit)
	if err != nil {
		slog.Error("select pvz failed", slog.Any("err", err))
		return nil, err
	}

	return pvzList, nil
}
