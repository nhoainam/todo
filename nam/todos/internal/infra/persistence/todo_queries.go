package persistence

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/gateway"
)

// Phase 2: TodoQueriesGateway GORM implementation. See resources/phase-02-database-di.md

type TodoQueriesGateway struct {
}

func NewTodoQueriesGateway() *TodoQueriesGateway {
	return &TodoQueriesGateway{}
}

func (g *TodoQueriesGateway) Get(context.Context, entity.TodoID, *gateway.GetTodoOptions) (*entity.Todo, error) {
	return nil, nil
}

func (g *TodoQueriesGateway) List(context.Context, *gateway.ListTodosOptions) (*gateway.OffsetPageResult[*entity.Todo], error) {
	return nil, nil
}
