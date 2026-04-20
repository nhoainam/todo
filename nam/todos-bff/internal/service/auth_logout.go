package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/output"
)

type authLogoutService struct {
	gateway gateway.AuthServiceGateway
}

func NewAuthLogoutService(gateway gateway.AuthServiceGateway) usecase.AuthLogout {
	return &authLogoutService{gateway: gateway}
}

func (s *authLogoutService) Logout(ctx context.Context, in *input.LogoutInput) (*output.LogoutOutput, error) {
	if err := s.gateway.Logout(ctx, in.UserID); err != nil {
		return nil, err
	}

	return &output.LogoutOutput{}, nil
}
