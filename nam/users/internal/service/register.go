package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase/output"
)

type userRegister struct {
	gateway gateway.AuthServiceGateway
}

func NewUserRegisterService(gateway gateway.AuthServiceGateway) usecase.UserRegister {
	return &userRegister{gateway: gateway}
}

func (s *userRegister) Register(ctx context.Context, input *input.RegisterInput) (*output.RegisterOutput, error) {
	user, err := s.gateway.Register(ctx, input.Username, input.Password)
	if err != nil {
		return nil, err
	}

	return &output.RegisterOutput{
		UserID:   user.ID.Int64(),
		Username: user.Username,
	}, nil
}
