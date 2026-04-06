package datastore

import (
	"context"
	"fmt"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/infra/query"
)

type todoWriter struct{}

func NewTodoWriter() gateway.TodoCommandsGateway {
	return &todoWriter{}
}

func (w *todoWriter) Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	db, err := DBFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get db from context: %w", err)
	}

	if err := query.Use(db).Todo.WithContext(ctx).Create(todo); err != nil {
		return nil, fmt.Errorf("create todo: %w", err)
	}

	return todo, nil
}

func (w *todoWriter) Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error) {
	db, err := DBFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get db from context: %w", err)
	}

	q := query.Use(db).Todo
	if _, err := q.WithContext(ctx).Where(q.ID.Eq(int64(todo.ID))).Updates(todo); err != nil {
		return nil, fmt.Errorf("update todo: %w", err)
	}

	return todo, nil
}

func (w *todoWriter) Delete(ctx context.Context, todoID entity.TodoID) error {
	db, err := DBFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get db from context: %w", err)
	}

	q := query.Use(db).Todo
	if _, err := q.WithContext(ctx).Where(q.ID.Eq(int64(todoID))).Delete(); err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}

	return nil
}
