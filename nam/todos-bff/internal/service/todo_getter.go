package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/output"
)

type todoService struct {
	todoGateway gateway.TodoServiceGateway
}

func NewTodoService(todoGateway gateway.TodoServiceGateway) usecase.TodoGetter {
	return &todoService{
		todoGateway: todoGateway,
	}
}

func (s *todoService) GetTodo(ctx context.Context, in *input.GetTodoInput) (*output.TodoOutput, error) {
	entityTodo, err := s.todoGateway.GetTodo(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &output.TodoOutput{
		Todo: entityTodo,
	}, nil
}

func (s *todoService) ListTodos(ctx context.Context, in *input.ListTodosInput) (*output.TodoListOutput, error) {
	result, err := s.todoGateway.ListTodos(ctx, in.Parent, &gateway.ListTodosOptions{
		Pagination: &gateway.OffsetPage{
			Limit:  in.Limit,
			Offset: in.Offset,
		},
	})
	if err != nil {
		return nil, err
	}

	return &output.TodoListOutput{
		Todos:      result.Items,
		ListName:   result.ListName,
		TotalCount: int(result.TotalCount),
	}, nil
}
