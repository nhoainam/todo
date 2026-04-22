package output

import "github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/entity"

type TodoOutput struct {
	*entity.Todo
}

type TodoListOutput struct {
	Todos      []*entity.Todo
	ListName   string
	TotalCount int
}

type CreateTodoOutput struct {
	*entity.Todo
}

type UpdateTodoOutput struct {
	*entity.Todo
}

type DeleteTodoOutput struct {
	// No fields needed for delete output, but you can add metadata if necessary
}
