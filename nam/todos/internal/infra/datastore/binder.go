package datastore

import (
	"context"

	"gorm.io/gorm"
)

type binder struct {
	db *gorm.DB
}

func NewBinder(db *gorm.DB) *binder {
	return &binder{db: db}
}

func (b *binder) Bind(ctx context.Context) context.Context {
	tx := b.db.Begin()
	return WithDB(ctx, tx)
}
