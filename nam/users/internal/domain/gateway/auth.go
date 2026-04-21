package gateway

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/domain/entity"
)

type AuthServiceGateway interface {
	Login(ctx context.Context, username, password string) (*entity.User, error)
	Logout(ctx context.Context, userID int64) error
	Register(ctx context.Context, username, password string) (*entity.User, error)
}
