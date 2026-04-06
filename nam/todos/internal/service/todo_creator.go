package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/output"
)

type todoCreator struct {
	todoCommandsGateway gateway.TodoCommandsGateway
	binder              gateway.Binder
}

// NewTodoCreator creates a new TodoCreator service.
func NewTodoCreator(c gateway.TodoCommandsGateway, b gateway.Binder) usecase.TodoCreator {
	return &todoCreator{
		todoCommandsGateway: c,
		binder:              b,
	}
}

func (s *todoCreator) Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error) {
	// 1. Start transaction
	ctx = s.binder.Bind(ctx)

	// 2. Always rollback on panic or early exit if not committed
	defer s.binder.Rollback(ctx)

	todo := &entity.Todo{
		ListID:    in.ListID,
		CreatorID: in.CreatorID,
		Title:     in.Title,
		Content:   in.Content,
		Status:    entity.TodoStatusPENDING,
		DueDate:   in.DueDate,
	}

	created, err := s.todoCommandsGateway.Create(ctx, todo)
	if err != nil {
		return nil, err
	}

	// 3. Commit if everything was successful
	s.binder.Commit(ctx)
	
	return &output.TodoCreator{Todo: created}, nil
}
