package gateway

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
)

// Phase 2: TodoCommandsGateway interface — write operations (Create, Update, Delete). See resources/phase-02-database-di.md

type TodoCommandsGateway interface {
	Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
	Delete(ctx context.Context, todoID entity.TodoID) error
	Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
}
