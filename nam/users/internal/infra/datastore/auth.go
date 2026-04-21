package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/domain/entity"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/infra/query"
	"gorm.io/gorm"
)

type userReader struct{}

func NewUserReader() gateway.AuthServiceGateway {
	return &userReader{}
}

func (r *userReader) Login(ctx context.Context, username, password string) (*entity.User, error) {
	db, err := DBFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get db from context: %w", err)
	}

	q := query.Use(db).User
	qb := q.WithContext(ctx)

	user, err := qb.Where(q.Username.Eq(username)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, fmt.Errorf("query user by username: %w", err)
	}

	if !verifyPassword(password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	if !isBcryptHash(user.Password) {
		hashedPassword, hashErr := hashPassword(password)
		if hashErr == nil {
			_, _ = qb.Where(q.ID.Eq(user.ID.Int64())).Update(q.Password, hashedPassword)
		}
	}

	return &entity.User{
		ID:       entity.UserID(user.ID),
		Username: user.Username,
	}, nil
}

func (r *userReader) Logout(ctx context.Context, userID int64) error {
	_, err := DBFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get db from context: %w", err)
	}

	return nil
}

func (r *userReader) Register(ctx context.Context, username, password string) (*entity.User, error) {
	db, err := DBFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get db from context: %w", err)
	}

	q := query.Use(db).User
	qb := q.WithContext(ctx)

	if _, err := qb.Where(q.Username.Eq(username)).First(); err == nil {
		return nil, errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("query existing user: %w", err)
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &entity.User{
		Username: username,
		Password: hashedPassword,
	}

	if err := qb.Create(newUser); err != nil {
		return nil, fmt.Errorf("create new user: %w", err)
	}

	return &entity.User{
		ID:       newUser.ID,
		Username: newUser.Username,
	}, nil
}
