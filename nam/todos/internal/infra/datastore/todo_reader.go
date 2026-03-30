// IMPLEMENTATION — Infrastructure layer
// File: todos/internal/infrastructure/datastore/todo_reader.go

package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/infra/persistence/mapper"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/infra/query"
	"gorm.io/gorm"
)

type todoReader struct{}

func NewTodoReader() gateway.TodoQueriesGateway {
	return &todoReader{}
}

func (r *todoReader) Get(
	ctx context.Context,
	todoID entity.TodoID,
	opts *gateway.GetTodoOptions,
) (*entity.Todo, error) {
	// 1. Get DB connection from context
	db, err := DBFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get db from context: %w", err)
	}

	q := query.Use(db).Todo
	qb := q.WithContext(ctx)

	// 2. Execute query
	todo, err := qb.Where(q.ID.Eq(int64(todoID))).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found -> return nil (NOT an error)
		}
		return nil, fmt.Errorf("get todo: %w", err)
	}

	return mapper.ToDomainTodo(todo), nil
}

// List returns a paginated list of todos, filtered by the options provided.
func (r *todoReader) List(
	ctx context.Context,
	opts *gateway.ListTodosOptions,
) (*gateway.OffsetPageResult[*entity.Todo], error) {
	db, err := DBFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get db from context: %w", err)
	}

	q := query.Use(db)
	qb := q.Todo.WithContext(ctx)

	// Apply optional filters
	if opts != nil {
		qb = r.applyFilter(q, qb, opts.Filter)
	}

	// Determine pagination parameters
	offset, limit := 0, 20 // sensible defaults
	if opts != nil && opts.Pagination != nil {
		offset = opts.Pagination.Offset
		limit = opts.Pagination.Limit
	}

	models, total, err := qb.FindByPage(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}

	todos := make([]*entity.Todo, 0, len(models))
	for _, m := range models {
		todos = append(todos, mapper.ToDomainTodo(m))
	}

	return &gateway.OffsetPageResult[*entity.Todo]{
		Items:      todos,
		TotalCount: total,
		Page:       &gateway.OffsetPage{Offset: offset, Limit: limit},
	}, nil
}

// applyFilter adds WHERE conditions to the query builder based on non-nil filter fields.
func (r *todoReader) applyFilter(q *query.Query, builder query.ITodoDo, filter *gateway.TodoFilter) query.ITodoDo {
	if filter == nil {
		return builder
	}

	todoQuery := q.Todo

	if filter.StatusEq != nil {
		statusMap := map[entity.TodoStatus]int{
			entity.TodoStatusPENDING:     0,
			entity.TodoStatusIN_PROGRESS: 1,
			entity.TodoStatusDONE:        2,
		}
		builder = builder.Where(todoQuery.Status.Eq(statusMap[*filter.StatusEq]))
	}

	if filter.PriorityEq != nil {
		builder = builder.Where(todoQuery.Priority.Eq(filter.PriorityEq.Int()))
	}

	if filter.DueDateGTE != nil {
		builder = builder.Where(todoQuery.DueDate.Gte(*filter.DueDateGTE))
	}

	if filter.DueDateLTE != nil {
		builder = builder.Where(todoQuery.DueDate.Lte(*filter.DueDateLTE))
	}

	return builder
}
