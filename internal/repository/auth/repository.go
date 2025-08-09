package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type repo struct {
	pool *pgxpool.Pool
}

func CreateRepository(ctx context.Context, dsn string) (*repo, error) {
	repository := &repo{}
	err := repository.openPool(ctx, dsn)
	return repository, err
}

func (rep *repo) openPool(ctx context.Context, dsn string) error {
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	rep.pool = pool
	return nil
}

func (rep *repo) ClosePool(ctx context.Context) {
	if rep.pool != nil {
		rep.pool.Close()
	}
}
