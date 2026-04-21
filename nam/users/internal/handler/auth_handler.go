package handler

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase/input"
	userv1 "github.com/tuannguyenandpadcojp/fresher26/nam/users/proto/user/v1"
)

type server struct {
	userv1.UnimplementedUserServiceServer
	login    usecase.UserLogin
	logout   usecase.UserLogout
	register usecase.UserRegister
}

func NewServer(login usecase.UserLogin, logout usecase.UserLogout, register usecase.UserRegister) userv1.UserServiceServer {
	return &server{
		login:    login,
		logout:   logout,
		register: register,
	}
}

func (s *server) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	input := &input.LoginInput{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	output, err := s.login.Login(ctx, input)
	if err != nil {
		return nil, err
	}

	return &userv1.LoginResponse{
		User: &userv1.User{
			Id:       output.UserID,
			Username: output.Username,
		},
		AccessToken: output.AccessToken,
	}, nil
}

func (s *server) Logout(ctx context.Context, req *userv1.LogoutRequest) (*userv1.LogoutResponse, error) {
	input := &input.LogoutInput{
		UserID: req.GetUserId(),
	}

	_, err := s.logout.Logout(ctx, input)
	if err != nil {
		return nil, err
	}

	return &userv1.LogoutResponse{}, nil
}

func (s *server) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	input := &input.RegisterInput{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	output, err := s.register.Register(ctx, input)
	if err != nil {
		return nil, err
	}

	return &userv1.RegisterResponse{
		User: &userv1.User{
			Id:       output.UserID,
			Username: output.Username,
		},
	}, nil
}
