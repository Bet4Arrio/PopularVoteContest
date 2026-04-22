package sqlrepo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(ctx context.Context, connString string) (*DB, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	db := &DB{Pool: pool}
	err = Migrate(ctx, db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
