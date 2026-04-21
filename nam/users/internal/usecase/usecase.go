package usecase

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase/output"
)

type UserLogin interface {
	Login(ctx context.Context, input *input.LoginInput) (*output.LoginOutput, error)
}

type UserLogout interface {
	Logout(ctx context.Context, input *input.LogoutInput) (*output.LogoutOutput, error)
}

type UserRegister interface {
	Register(ctx context.Context, input *input.RegisterInput) (*output.RegisterOutput, error)
}
type TokenGenerator interface {
	GenerateToken(userID int64, username string) (string, error)
}
