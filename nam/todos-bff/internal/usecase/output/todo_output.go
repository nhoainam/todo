package output

import "github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/entity"

type TodoOutput struct {
	*entity.Todo
}

type TodoListOutput struct {
	Todos      []*entity.Todo
	Listname   string
	TotalCount int
}
