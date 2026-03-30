package service

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/infra/idgen"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/output"
)

type todoCreator struct {
	todoCommandsGateway gateway.TodoCommandsGateway
	idGen               idgen.IDGenerator
}

// NewTodoCreator creates a new TodoCreator service.
func NewTodoCreator(c gateway.TodoCommandsGateway, idGen idgen.IDGenerator) usecase.TodoCreator {
	return &todoCreator{
		todoCommandsGateway: c,
		idGen:               idGen,
	}
}

func (s *todoCreator) Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error) {
	todo := &entity.Todo{
		ID:        s.idGen.NewTodoID(),
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

	return &output.TodoCreator{Todo: created}, nil
}
