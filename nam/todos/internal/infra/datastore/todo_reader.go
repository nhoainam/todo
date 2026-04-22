// IMPLEMENTATION — Infrastructure layer
// File: todos/internal/infrastructure/datastore/todo_reader.go

package datastore

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/apperrors"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/gateway"
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
	qb := q.WithContext(ctx).Where(q.ID.Eq(int64(todoID)))

	if opts != nil {
		if opts.ListIDEq != nil {
			qb = qb.Where(q.ListID.Eq(int64(*opts.ListIDEq)))
		}
		if opts.CreatorIDEq != nil {
			qb = qb.Where(q.CreatorID.Eq(int64(*opts.CreatorIDEq)))
		}
	}

	// 2. Execute query
	todo, err := qb.First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found -> return nil (NOT an error)
		}
		return nil, fmt.Errorf("get todo: %w", err)
	}

	return todo, nil
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

	// Determine pagination parameters
	offset, limit := 0, 20 // sensible defaults
	if opts != nil && opts.Pagination != nil {
		offset = opts.Pagination.Offset
		limit = opts.Pagination.Limit
	}

	var listName string
	// If there's a parent list filter, validate it and get the list name for the response

	if opts != nil {
		exists, parentListName, err := r.ensureParentListOwnedByUser(ctx, q, opts.Filter)
		if err != nil {
			return nil, fmt.Errorf("validate parent todo list: %w", err)
		}
		if parentListName != nil {
			listName = *parentListName
		}
		if !exists {
			return nil, apperrors.NewNotFound("todo list not found", nil)
		}
	}

	qb := q.Todo.WithContext(ctx)

	// Apply optional filters
	if opts != nil {
		qb = r.applyFilter(q, qb, opts.Filter)
	}

	models, total, err := qb.FindByPage(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}

	todos := make([]*entity.Todo, 0, len(models))
	for _, m := range models {
		todos = append(todos, m)
	}

	return &gateway.OffsetPageResult[*entity.Todo]{
		Items:      todos,
		TotalCount: total,
		ListName:   listName,
		Page:       &gateway.OffsetPage{Offset: offset, Limit: limit},
	}, nil
}

// ensureParentListOwnedByUser validates the parent todo list before querying child todos.
func (r *todoReader) ensureParentListOwnedByUser(
	ctx context.Context,
	q *query.Query,
	filter *gateway.TodoFilter,
) (bool, *string, error) {
	if filter == nil || filter.ListIDEq == nil || filter.CreatorIDEq == nil {
		return true, nil, nil // No parent list filter → skip validation
	}

	todoList, err := q.TodoList.WithContext(ctx).Where(
		q.TodoList.ID.Eq(int64(*filter.ListIDEq)),
		q.TodoList.OwnerId.Eq(int64(*filter.CreatorIDEq)),
	).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		return false, nil, fmt.Errorf("get todo list: %w", err)
	}

	return true, &todoList.Name, nil
}

// applyFilter adds WHERE conditions to the query builder based on non-nil filter fields.
func (r *todoReader) applyFilter(q *query.Query, builder query.ITodoDo, filter *gateway.TodoFilter) query.ITodoDo {
	if filter == nil {
		return builder
	}

	todoQuery := q.Todo

	if filter.ListIDEq != nil {
		builder = builder.Where(todoQuery.ListID.Eq(int64(*filter.ListIDEq)))
	}

	if filter.CreatorIDEq != nil {
		builder = builder.Where(todoQuery.CreatorID.Eq(int64(*filter.CreatorIDEq)))
	}

	if filter.StatusEq != nil {
		statusCode, err := filter.StatusEq.DBCode()
		if err == nil {
			statusCodeText := strconv.FormatInt(statusCode, 10)
			if *filter.StatusEq == entity.TodoStatusPENDING {
				builder = builder.Where(todoQuery.Status.In("0", statusCodeText))
			} else {
				builder = builder.Where(todoQuery.Status.Eq(statusCodeText))
			}
		}
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
