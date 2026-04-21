package usecase

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/output"
)

type TodoGetter interface {
	GetTodo(ctx context.Context, in *input.GetTodoInput) (*output.TodoOutput, error)
	ListTodos(ctx context.Context, in *input.ListTodosInput) (*output.TodoListOutput, error)
}

type TodoUpdater interface {
	UpdateTodo(ctx context.Context, in *input.UpdateTodoInput) (*output.TodoOutput, error)
}

type AuthLogin interface {
	Login(ctx context.Context, in *input.LoginInput) (*output.LoginOutput, error)
}

type AuthLogout interface {
	Logout(ctx context.Context, in *input.LogoutInput) (*output.LogoutOutput, error)
}

type TokenGenerator interface {
	GenerateToken(userID int64, username string) (string, error)
}

type AuthRegister interface {
	Register(ctx context.Context, in *input.RegisterInput) (*output.RegisterOutput, error)
}
