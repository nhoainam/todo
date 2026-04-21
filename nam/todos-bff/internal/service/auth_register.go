package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/apperrors"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/output"
)

type authRegisterService struct {
	gateway gateway.AuthServiceGateway
}

func NewAuthRegisterService(gateway gateway.AuthServiceGateway) usecase.AuthRegister {
	return &authRegisterService{gateway: gateway}
}

func (s *authRegisterService) Register(ctx context.Context, in *input.RegisterInput) (*output.RegisterOutput, error) {
	user, err := s.gateway.Register(ctx, in.Username, in.Password)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.NewInternal("register response user is nil", nil)
	}

	return &output.RegisterOutput{User: user}, nil
}
