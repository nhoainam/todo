package service

import (
	"context"
	"fmt"

	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase/output"
)

type LoginService struct {
	gateway        gateway.AuthServiceGateway
	tokenGenerator usecase.TokenGenerator
}

func NewLoginService(gateway gateway.AuthServiceGateway, tokenGenerator usecase.TokenGenerator) usecase.UserLogin {
	return &LoginService{
		gateway:        gateway,
		tokenGenerator: tokenGenerator,
	}
}

func (s *LoginService) Login(ctx context.Context, input *input.LoginInput) (*output.LoginOutput, error) {
	user, err := s.gateway.Login(ctx, input.Username, input.Password)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenGenerator.GenerateToken(user.ID.Int64(), user.Username)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &output.LoginOutput{
		UserID:      user.ID.Int64(),
		Username:    user.Username,
		AccessToken: accessToken,
	}, nil
}
