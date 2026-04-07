package dataloader

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/handler/graph/model"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/handler/graph/scalar"
)

type Loaders struct {
	UserLoader     *dataloader.Loader[string, *model.User]
	TodoListLoader *dataloader.Loader[string, *model.TodoList]
}

type userBatchLoader struct{}

func (b *userBatchLoader) batchGetUsers(ctx context.Context, keys []string) []*dataloader.Result[*model.User] {
	results := make([]*dataloader.Result[*model.User], len(keys))
	for i, key := range keys {
		results[i] = &dataloader.Result[*model.User]{
			Data: &model.User{Name: scalar.ResourceName(key)},
			Error: nil,
		}
	}
	return results
}

type todoListBatchLoader struct{}

func (b *todoListBatchLoader) batchGetTodoLists(ctx context.Context, keys []string) []*dataloader.Result[*model.TodoList] {
	results := make([]*dataloader.Result[*model.TodoList], len(keys))
	for i, key := range keys {
		results[i] = &dataloader.Result[*model.TodoList]{
			Data: &model.TodoList{Name: scalar.ResourceName(key)},
			Error: nil,
		}
	}
	return results
}

func NewLoaders() *Loaders {
	userLoader := &userBatchLoader{}
	todoListLoader := &todoListBatchLoader{}

	return &Loaders{
		UserLoader: dataloader.NewBatchedLoader(
			userLoader.batchGetUsers,
			dataloader.WithCache[string, *model.User](&dataloader.NoCache[string, *model.User]{}),
			dataloader.WithWait[string, *model.User](time.Millisecond*5),
		),
		TodoListLoader: dataloader.NewBatchedLoader(
			todoListLoader.batchGetTodoLists,
			dataloader.WithCache[string, *model.TodoList](&dataloader.NoCache[string, *model.TodoList]{}),
			dataloader.WithWait[string, *model.TodoList](time.Millisecond*5),
		),
	}
}

type key string

const loadersKey key = "dataloaders"

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

func WithLoaders(ctx context.Context, loaders *Loaders) context.Context {
	return context.WithValue(ctx, loadersKey, loaders)
}
