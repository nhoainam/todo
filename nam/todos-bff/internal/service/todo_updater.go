package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/output"
)

type TodoUpdater struct {
	todoGateway gateway.TodoServiceGateway
}

func NewTodoUpdater(todoGateway gateway.TodoServiceGateway) usecase.TodoUpdater {
	return &TodoUpdater{
		todoGateway: todoGateway,
	}
}

func (s *TodoUpdater) UpdateTodo(ctx context.Context, in *input.UpdateTodoInput) (*output.TodoOutput, error) {
	entityTodo, err := s.todoGateway.UpdateTodo(ctx, &gateway.UpdateTodoInput{
		Name:    in.Name,
		Title:   in.Title,
		Content: in.Content,
		Status:  in.Status,
		DueDate: in.DueDate,
	})
	if err != nil {
		return nil, err
	}

	return &output.TodoOutput{Todo: entityTodo}, nil
}
