package grpc_client

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/apperrors"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/entity"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/gateway"
	userv1 "github.com/tuannguyenandpadcojp/fresher26/nam/users/proto/user/v1"
)

type authServiceClient struct {
	client userv1.UserServiceClient
}

func NewAuthServiceClient(client userv1.UserServiceClient) gateway.AuthServiceGateway {
	return &authServiceClient{client: client}
}

func (c *authServiceClient) Login(ctx context.Context, username, password string) (*entity.User, string, error) {
	resp, err := c.client.Login(ctx, &userv1.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, "", mapGRPCError(err)
	}

	user := resp.GetUser()
	if user == nil {
		return nil, "", apperrors.NewInternal("login response user is nil", nil)
	}

	return &entity.User{
		ID:       entity.UserID(user.GetId()),
		Username: user.GetUsername(),
	}, resp.GetAccessToken(), nil
}

func (c *authServiceClient) Logout(ctx context.Context, userID int64) error {
	_, err := c.client.Logout(ctx, &userv1.LogoutRequest{UserId: userID})
	if err != nil {
		return mapGRPCError(err)
	}

	return nil
}

func (c *authServiceClient) Register(ctx context.Context, username, password string) (*entity.User, error) {
	resp, err := c.client.Register(ctx, &userv1.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, mapGRPCError(err)
	}

	user := resp.GetUser()
	if user == nil {
		return nil, apperrors.NewInternal("register response user is nil", nil)
	}

	return &entity.User{
		ID:       entity.UserID(user.GetId()),
		Username: user.GetUsername(),
	}, nil
}
