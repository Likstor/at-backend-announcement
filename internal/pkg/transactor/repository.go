package transactor

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryWithTransactor interface {
	Conn(ctx context.Context) Executor
}

type repositoryWithTransactor struct {
	pool *pgxpool.Pool
}

func NewRepositoryWithTransactor(pool *pgxpool.Pool) *repositoryWithTransactor {
	return &repositoryWithTransactor{
		pool: pool,
	}
}

func (repo repositoryWithTransactor) Conn(ctx context.Context) Executor {
	if tx := ExtractTx(ctx); tx != nil {
		return tx
	}

	return repo.pool
}