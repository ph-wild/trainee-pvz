package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectDB(ctx context.Context, connString string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
