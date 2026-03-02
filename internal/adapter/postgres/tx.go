package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitOfWork struct {
	UserPhotoRepo *UserPhotoRepo
}

func (m *TxManager) newUOW(db DBTX) *UnitOfWork {
	return &UnitOfWork{
		UserPhotoRepo: NewUserPhotoRepo(db),
	}
}

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

func (m *TxManager) WithTx(
	ctx context.Context,
	fn func(uow *UnitOfWork) error,
) error {

	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	uow := m.newUOW(tx)

	if err := fn(uow); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
