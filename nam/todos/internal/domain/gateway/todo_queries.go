package gateway

import (
	"context"
	"time"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
)

// Phase 2: TodoQueriesGateway interface — read operations (Get, List). See resources/phase-02-database-di.md

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
	ListIDEq    *entity.TodoListID
	CreatorIDEq *entity.UserID
	StatusEq    *entity.TodoStatus
	PriorityEq  *entity.Priority
	DueDateGTE  *time.Time
	DueDateLTE  *time.Time
}

type GetTodoOptions struct {
	ListIDEq    *entity.TodoListID
	CreatorIDEq *entity.UserID
}

type ListTodosOptions struct {
	Filter     *TodoFilter
	Pagination *OffsetPage
}

type TodoQueriesGateway interface {
	Get(ctx context.Context, todoID entity.TodoID, opts *GetTodoOptions) (*entity.Todo, error)
	List(ctx context.Context, opts *ListTodosOptions) (*OffsetPageResult[*entity.Todo], error)
}
