package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	er "trainee-pvz/internal/errors"
	"trainee-pvz/internal/models"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Add(ctx context.Context, p models.Product) error {
	query := `INSERT INTO products (id, datetime, type, reception_id) VALUES (:id, :datetime, :type, :reception_id)`
	_, err := r.db.NamedExecContext(ctx, query, p)
	if err != nil {
		slog.Error("add product failed", slog.Any("err", err))
		return errors.Wrap(err, "product repo: add product")
	}
	return nil
}

func (r *ProductRepository) DeleteLast(ctx context.Context, pvzID string) error {
	var receptionID string
	queryReception := `SELECT id FROM receptions WHERE pvz_id = $1 AND status = 'in_progress' ORDER BY datetime DESC LIMIT 1`
	err := r.db.GetContext(ctx, &receptionID, queryReception, pvzID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return er.ErrNoProducts
		}
		slog.Error("no product found", slog.Any("err", err))
		return errors.Wrap(err, "get last product id")
	}
	var productID string
	queryProduct := `
		SELECT id FROM products
		WHERE reception_id = $1
		ORDER BY datetime DESC
		LIMIT 1
	`
	err = r.db.GetContext(ctx, &productID, queryProduct, receptionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return er.ErrNoProducts
		}
		return errors.Wrap(err, "get last product id")
	}

	_, err = r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, productID)
	if err != nil {
		slog.Error("can't delete product", slog.Any("err", err))
		return errors.Wrap(err, "delete product")
	}
	return nil
}
