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

func (r *PVZRepository) ListWithReceptionsAndProducts(ctx context.Context, start, end *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	offset := (page - 1) * limit
	var pvzList []models.PVZ

	pvzQuery := `SELECT id, city, registration_date FROM pvz ORDER BY registration_date DESC OFFSET $1 LIMIT $2`
	err := r.db.SelectContext(ctx, &pvzList, pvzQuery, offset, limit)
	if err != nil {
		slog.Error("select pvz failed", slog.Any("err", err))
		return nil, err
	}

	var result []models.PVZWithReceptions
	for _, pvz := range pvzList {
		var receptions []models.Reception
		recQuery := `
			SELECT id, datetime, pvz_id, status
			FROM receptions
			WHERE pvz_id = $1
			AND ($2::timestamptz IS NULL OR datetime >= $2)
			AND ($3::timestamptz IS NULL OR datetime <= $3)`

		err := r.db.SelectContext(ctx, &receptions, recQuery, pvz.ID, start, end)
		if err != nil {
			slog.Error("select receptions failed", slog.Any("err", err))
			return nil, err
		}

		var recsWithProducts []models.ReceptionWithProducts
		for _, reception := range receptions {
			var products []models.Product

			prodQuery := `SELECT id, datetime, type, reception_id FROM products WHERE reception_id = $1 ORDER BY datetime DESC`
			err = r.db.SelectContext(ctx, &products, prodQuery, reception.ID)
			if err != nil {
				slog.Error("select products failed", slog.Any("err", err))
				return nil, err
			}

			recsWithProducts = append(recsWithProducts, models.ReceptionWithProducts{
				Reception: reception,
				Products:  products,
			})
		}

		result = append(result, models.PVZWithReceptions{
			PVZ:        pvz,
			Receptions: recsWithProducts,
		})
	}

	return result, nil
}
