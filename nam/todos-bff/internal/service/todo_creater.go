package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/output"
)

type TodoCreater struct {
	todoGateway gateway.TodoServiceGateway
}

func NewTodoCreater(todoGateway gateway.TodoServiceGateway) usecase.TodoCreator {
	return &TodoCreater{
		todoGateway: todoGateway,
	}
}

func (s *TodoCreater) Create(ctx context.Context, input *input.CreateTodoInput) (*output.CreateTodoOutput, error) {
	entityTodo, err := s.todoGateway.CreateTodo(ctx, input.Parent, &gateway.CreateTodoInput{
		Title:   input.Title,
		Content: input.Content,
		Status:  input.Status,
		DueDate: input.DueDate,
	})
	if err != nil {
		return nil, err
	}

	return &output.CreateTodoOutput{
		Todo: entityTodo,
	}, nil
}
