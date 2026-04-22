package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/output"
)

type todoLister struct {
	todoQueriesGateway gateway.TodoQueriesGateway
}

// NewTodoLister creates a new TodoLister service.
func NewTodoLister(q gateway.TodoQueriesGateway) usecase.TodoLister {
	return &todoLister{todoQueriesGateway: q}
}

func (l *todoLister) List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error) {
	listID := in.Name.TodoListID
	userID := in.Name.UserID
	opts := &gateway.ListTodosOptions{
		Filter: &gateway.TodoFilter{
			ListIDEq:    &listID,
			CreatorIDEq: &userID,
			StatusEq:    in.Status,
		},
		Pagination: &gateway.OffsetPage{
			Offset: int(in.Offset),
			Limit:  int(in.Limit),
		},
	}

	result, err := l.todoQueriesGateway.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &output.TodoLister{
		Todos:      result.Items,
		ListName:   result.ListName,
		TotalCount: int32(result.TotalCount),
	}, nil
}
