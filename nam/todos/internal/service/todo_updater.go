package service

import (
	"context"
	"time"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/apperrors"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/output"
)

type todoUpdater struct {
	todoQueriesGateway  gateway.TodoQueriesGateway
	todoCommandsGateway gateway.TodoCommandsGateway
}

// NewTodoUpdater creates a new TodoUpdater service.
func NewTodoUpdater(q gateway.TodoQueriesGateway, c gateway.TodoCommandsGateway) usecase.TodoUpdater {
	return &todoUpdater{
		todoQueriesGateway:  q,
		todoCommandsGateway: c,
	}
}

// Update implements the Read-Modify-Write pattern:
// 1. Fetch the existing todo to ensure it exists.
// 2. Apply only the fields present in the input DTO.
// 3. Persist the changes.
func (u *todoUpdater) Update(ctx context.Context, in *input.TodoUpdater) (*output.TodoUpdater, error) {
	existing, err := u.todoQueriesGateway.Get(ctx, in.Name.TodoID, nil)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, apperrors.NewNotFound("todo not found", nil)
	}

	if in.Title != nil {
		existing.Title = *in.Title
	}
	if in.Content != nil {
		existing.Content = *in.Content
	}
	if in.Status != nil {
		existing.Status = *in.Status
	}
	if in.DueDate != nil {
		existing.DueDate = in.DueDate
	}
	existing.UpdatedAt = time.Now()

	updated, err := u.todoCommandsGateway.Update(ctx, existing)
	if err != nil {
		return nil, err
	}
	return &output.TodoUpdater{Todo: updated}, nil
}
