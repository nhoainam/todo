package usecase

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/output"
)

type TodoGetter interface {
	Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error)
}
type TodoCreator interface {
	Create(ctx context.Context, input *input.TodoCreator) (*output.TodoCreator, error)
}
type TodoUpdater interface {
	Update(ctx context.Context, in *input.TodoUpdater) (*output.TodoUpdater, error)
}
type TodoLister interface {
	List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error)
}
type TodoDeleter interface {
	Delete(ctx context.Context, in *input.TodoDeleter) error
}
