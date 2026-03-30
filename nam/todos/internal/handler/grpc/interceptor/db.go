package interceptor

import (
	"context"

	"google.golang.org/grpc"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/infra/datastore"
	"gorm.io/gorm"
)

// DBInterceptor is a gRPC unary interceptor that injects a *gorm.DB into
// every request context. This allows gateway implementations to retrieve
// the DB via datastore.DBFromContext(ctx) without being passed it explicitly.
type DBInterceptor struct {
	db *gorm.DB
}

// NewDBInterceptor creates a new DBInterceptor provider.
func NewDBInterceptor(db *gorm.DB) *DBInterceptor {
	return &DBInterceptor{db: db}
}

// Unary returns a gRPC UnaryServerInterceptor that binds the DB to the context.
// If the DBInterceptor was created with a nil db, the interceptor is a no-op.
func (i *DBInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if i.db != nil {
			ctx = datastore.WithDB(ctx, i.db)
		}
		return handler(ctx, req)
	}
}
