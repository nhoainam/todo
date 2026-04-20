package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase/output"
)

type userLogout struct {
	gateway gateway.AuthServiceGateway
}

func NewUserLogoutService(gateway gateway.AuthServiceGateway) usecase.UserLogout {
	return &userLogout{gateway: gateway}
}

func (s *userLogout) Logout(ctx context.Context, input *input.LogoutInput) (*output.LogoutOutput, error) {
	if err := s.gateway.Logout(ctx, input.UserID); err != nil {
		return nil, err
	}

	return &output.LogoutOutput{}, nil
}
