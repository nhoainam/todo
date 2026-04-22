package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/apperrors"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/output"
)

type todoGetter struct {
	todoQueriesGateway gateway.TodoQueriesGateway
}

// NewTodoGetter creates a new instance of TodoGetter.
func NewTodoGetter(todoQueriesGateway gateway.TodoQueriesGateway) usecase.TodoGetter {
	return &todoGetter{
		todoQueriesGateway: todoQueriesGateway,
	}
}

func (g *todoGetter) Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error) {
	listID := in.Name.TodoListID
	userID := in.Name.UserID
	todo, err := g.todoQueriesGateway.Get(ctx, in.Name.TodoID, &gateway.GetTodoOptions{
		ListIDEq:    &listID,
		CreatorIDEq: &userID,
	})
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, apperrors.NewNotFound("todo not found", nil)
	}
	return &output.TodoGetter{Todo: todo}, nil
}
