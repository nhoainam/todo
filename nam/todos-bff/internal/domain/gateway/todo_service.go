package gateway

import (
	"context"
	"time"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/entity"
)

// todo_service.go — TodoService Gateway Interface
//
// Phase 3: GraphQL & BFF Pattern
//
// This file defines the interface for communicating with the backend Todos gRPC service.
//
//	type TodoServiceGateway interface {
//	    GetTodo(ctx context.Context, name string) (*domain.Todo, error)
//	    ListTodos(ctx context.Context, listName string, opts *ListTodosOptions) ([]*domain.Todo, error)
//	    CreateTodo(ctx context.Context, input *input.CreateTodoInput) (*domain.Todo, error)
//	    UpdateTodo(ctx context.Context, input *input.UpdateTodoInput) (*domain.Todo, error)
//	    DeleteTodo(ctx context.Context, name string) error
//	}
//
// The implementation (in infra/grpc_client/) translates these calls
// to actual gRPC requests to the backend service.
//
// See: resources/phase-03-graphql-bff.md (gateway pattern in BFF)
// OffsetPage holds pagination parameters for offset-based queries.
type OffsetPage struct {
	Offset int
	Limit  int
}

// OffsetPageResult is a generic paginated result set.
type OffsetPageResult[T any] struct {
	Items      []T
	TotalCount int64
	ListName   string
	Page       *OffsetPage
}

// TodoFilter holds optional filter conditions for listing todos.
// All fields are pointers — nil means "no filter on this field".
type TodoFilter struct {
	StatusEq   *entity.TodoStatus
	PriorityEq *entity.Priority
	DueDateGTE *time.Time
	DueDateLTE *time.Time
}

type GetTodoOptions struct {
	// Add any options for the Get operation
}

type ListTodosOptions struct {
	Filter     *TodoFilter
	Pagination *OffsetPage
}

type UpdateTodoInput struct {
	Name    string
	Title   *string
	Content *string
	Status  *entity.TodoStatus
	DueDate *time.Time
}

type CreateTodoInput struct {
	Title   string
	Content *string
	Status  entity.TodoStatus
	DueDate *time.Time
}

type DeleteTodoInput struct {
	Name string
}

type TodoServiceGateway interface {
	GetTodo(ctx context.Context, name string) (*entity.Todo, error)
	ListTodos(ctx context.Context, parent string, opts *ListTodosOptions) (*OffsetPageResult[*entity.Todo], error)
	CreateTodo(ctx context.Context, parent string, input *CreateTodoInput) (*entity.Todo, error)
	UpdateTodo(ctx context.Context, input *UpdateTodoInput) (*entity.Todo, error)
	DeleteTodo(ctx context.Context, input *DeleteTodoInput) error
}
