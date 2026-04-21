package datastore

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type dbKey struct{}

func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey{}, db)
}

func DBFromContext(ctx context.Context) (*gorm.DB, error) {
	db, ok := ctx.Value(dbKey{}).(*gorm.DB)
	if !ok {
		return nil, errors.New("db not found in context")
	}
	return db, nil
}
