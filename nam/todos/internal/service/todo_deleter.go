package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/apperrors"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/input"
)

type todoDeleter struct {
	todoQueriesGateway  gateway.TodoQueriesGateway
	todoCommandsGateway gateway.TodoCommandsGateway
}

// NewTodoDeleter creates a new TodoDeleter service.
func NewTodoDeleter(q gateway.TodoQueriesGateway, c gateway.TodoCommandsGateway) usecase.TodoDeleter {
	return &todoDeleter{
		todoQueriesGateway:  q,
		todoCommandsGateway: c,
	}
}

func (s *todoDeleter) Delete(ctx context.Context, in *input.TodoDeleter) error {
	listID := in.Name.TodoListID
	userID := in.Name.UserID
	existing, err := s.todoQueriesGateway.Get(ctx, in.Name.TodoID, &gateway.GetTodoOptions{
		ListIDEq:    &listID,
		CreatorIDEq: &userID,
	})
	if err != nil {
		return err
	}
	if existing == nil {
		return apperrors.NewNotFound("todo not found", nil)
	}

	return s.todoCommandsGateway.Delete(ctx, in.Name.TodoID)
}
