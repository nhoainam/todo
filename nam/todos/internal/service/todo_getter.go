package service

import (
	"context"

	apperrors "github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain"
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
	todo, err := g.todoQueriesGateway.Get(ctx, in.Name.TodoID, nil)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, apperrors.NewNotFound("todo not found")
	}
	return &output.TodoGetter{Todo: todo}, nil
}
