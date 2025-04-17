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

type ReceptionRepository struct {
	db *sqlx.DB
}

func NewReceptionRepository(db *sqlx.DB) *ReceptionRepository {
	return &ReceptionRepository{db: db}
}

func (r *ReceptionRepository) HasOpenReception(ctx context.Context, pvzID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM receptions WHERE pvz_id = $1 AND status = 'in_progress'`
	err := r.db.GetContext(ctx, &count, query, pvzID)
	return count > 0, err
}
func (r *ReceptionRepository) GetLastReceptionID(ctx context.Context, pvzID string) (string, error) {
	var id string
	query := `
		SELECT id FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY datetime DESC
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &id, query, pvzID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", er.ErrNoOpenReception // тут ли?
		}
		slog.Error("failed to get last reception", slog.Any("err", err))
		return "", errors.Wrap(err, "reception repo: cget last reception")
	}
	return id, nil
}

func (r *ReceptionRepository) Create(ctx context.Context, rec models.Reception) error {
	query := `INSERT INTO receptions (id, datetime, pvz_id, status) VALUES (:id, :datetime, :pvz_id, :status)`
	_, err := r.db.NamedExecContext(ctx, query, rec)
	if err != nil {
		slog.Error("create reception failed", slog.Any("err", err))
		return errors.Wrap(err, "reception repo: create reception")
	}
	return nil
}

func (r *ReceptionRepository) Close(ctx context.Context, id string) error {
	query := `UPDATE receptions SET status = 'close' WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("close reception failed", slog.Any("err", err))
		return errors.Wrap(err, "reception repo: close reception")
	}
	return nil
}

func (r *ReceptionRepository) GetOpenReceptionID(ctx context.Context, pvzID string) (string, error) {
	var id string
	query := `
		SELECT id FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY datetime DESC
		LIMIT 1
	`
	//TODO in_progress
	err := r.db.GetContext(ctx, &id, query, pvzID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", er.ErrNoOpenReception
		}
		slog.Error("get open reception failed", slog.Any("err", err))
		return "", errors.Wrap(err, "get open reception failed")
	}

	return id, nil
}
