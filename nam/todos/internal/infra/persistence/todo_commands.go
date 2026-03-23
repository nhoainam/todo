package persistence

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
)

// Phase 2: TodoCommandsGateway GORM implementation. See resources/phase-02-database-di.md

type TodoCommandsGateway struct {
}

func NewTodoCommandsGateway() *TodoCommandsGateway {
	return &TodoCommandsGateway{}
}

func (g *TodoCommandsGateway) Create(context.Context, *entity.Todo) (*entity.Todo, error) {
	return nil, nil
}
func (g *TodoCommandsGateway) Update(context.Context, *entity.Todo) (*entity.Todo, error) {
	return nil, nil
}
func (g *TodoCommandsGateway) Delete(context.Context, entity.TodoID) error { return nil }
