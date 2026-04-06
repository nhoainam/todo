package datastore

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/gateway"
	"gorm.io/gorm"
)

type binder struct {
	db *gorm.DB
}

// NewBinder returns a Binder implementation handling gorm DB transactions.
func NewBinder(db *gorm.DB) gateway.Binder {
	return &binder{db: db}
}

func (b *binder) Bind(ctx context.Context) context.Context {
	tx := b.db.WithContext(ctx).Begin()
	return WithDB(ctx, tx)
}

func (b *binder) Commit(ctx context.Context) error {
	db, err := DBFromContext(ctx)
	if err != nil {
		return err
	}
	return db.Commit().Error
}

func (b *binder) Rollback(ctx context.Context) error {
	db, err := DBFromContext(ctx)
	if err != nil {
		return err
	}
	return db.Rollback().Error
}
