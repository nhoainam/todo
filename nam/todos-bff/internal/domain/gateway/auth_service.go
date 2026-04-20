package gateway

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/entity"
)

type AuthServiceGateway interface {
	Login(ctx context.Context, username, password string) (*entity.User, string, error)
	Logout(ctx context.Context, userID int64) error
	Register(ctx context.Context, username, password string) (*entity.User, error)
}
