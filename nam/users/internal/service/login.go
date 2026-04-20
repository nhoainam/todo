package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase/output"
)

type LoginService struct {
	gateway gateway.AuthServiceGateway
}

func NewLoginService(gateway gateway.AuthServiceGateway) usecase.UserLogin {
	return &LoginService{
		gateway: gateway,
	}
}

func (s *LoginService) Login(ctx context.Context, input *input.LoginInput) (*output.LoginOutput, error) {
	user, err := s.gateway.Login(ctx, input.Username, input.Password)
	if err != nil {
		return nil, err
	}

	accessToken, err := generateToken(user.ID.Int64(), user.Username)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &output.LoginOutput{
		UserID:      user.ID.Int64(),
		Username:    user.Username,
		AccessToken: accessToken,
	}, nil
}

func generateToken(userID int64, username string) (string, error) {
	if userID <= 0 {
		return "", errors.New("invalid user id")
	}
	if strings.TrimSpace(username) == "" {
		return "", errors.New("username is required")
	}

	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	now := time.Now().UTC()

	claims := jwt.MapClaims{
		"sub":      strconv.FormatInt(userID, 10),
		"user_id":  userID,
		"username": username,
		"iat":      now.Unix(),
		"exp":      now.Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign jwt payload: %w", err)
	}

	return signedToken, nil
}
