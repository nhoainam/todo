package input

import (
	"time"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/entity"
)

type GetTodoInput struct {
	Name string
}

type CreateTodoInput struct {
	Title       string
	Description string
}

type UpdateTodoInput struct {
	Name    string
	Title   *string
	Content *string
	Status  *entity.TodoStatus
	DueDate *time.Time
}

type ListTodosInput struct {
	Parent string
	Limit  int
	Offset int
}
