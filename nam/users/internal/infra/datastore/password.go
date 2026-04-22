package datastore

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(plainPassword string) (string, error) {
	if strings.TrimSpace(plainPassword) == "" {
		return "", errors.New("password is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate bcrypt hash: %w", err)
	}

	return string(hashedPassword), nil
}

func verifyPassword(plainPassword, storedPassword string) bool {
	if storedPassword == "" {
		return false
	}

	if isBcryptHash(storedPassword) {
		return bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(plainPassword)) == nil
	}

	return plainPassword == storedPassword
}

func isBcryptHash(value string) bool {
	return strings.HasPrefix(value, "$2a$") || strings.HasPrefix(value, "$2b$") || strings.HasPrefix(value, "$2y$")
}
