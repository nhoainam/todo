package usecase

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/output"
)

type TodoGetter interface {
	GetTodo(ctx context.Context, in *input.GetTodoInput) (*output.TodoOutput, error)
}
