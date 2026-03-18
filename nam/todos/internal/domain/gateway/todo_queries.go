package gateway

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
)

// Phase 2: TodoQueriesGateway interface — read operations (Get, List). See resources/phase-02-database-di.md

type GetTodoOptions struct {
	// Add any options for the Get operation
}

type ListTodosOptions struct {
	// Add any options for the List operation
}

type TodoQueriesGateway interface {
	Get(ctx context.Context, todoID entity.TodoID, opts *GetTodoOptions) (*entity.Todo, error)
	// List(ctx context.Context, opts *ListTodosOptions) (*query.OffsetPageResult[*entity.Todo], error)
}
