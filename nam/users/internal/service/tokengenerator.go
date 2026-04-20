package service

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/usecase"
)

type tokenGeneratorService struct{}

func NewTokenGeneratorService() usecase.TokenGenerator {
	return &tokenGeneratorService{}
}

func (s *tokenGeneratorService) GenerateToken(userID int64, username string) (string, error) {
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
