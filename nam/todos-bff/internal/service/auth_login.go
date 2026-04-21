package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/apperrors"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/output"
)

func NewAuthLoginService(gateway gateway.AuthServiceGateway) usecase.AuthLogin {
	return &authLoginService{gateway: gateway}
}

type authLoginService struct {
	gateway gateway.AuthServiceGateway
}

func (s *authLoginService) Login(ctx context.Context, in *input.LoginInput) (*output.LoginOutput, error) {
	user, accessToken, err := s.gateway.Login(ctx, in.Username, in.Password)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.NewInternal("login response user is nil", nil)
	}

	return &output.LoginOutput{
		User:        user,
		AccessToken: accessToken,
	}, nil
}
