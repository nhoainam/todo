package mapper

import (
	"fmt"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/apperrors"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/entity"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/handler/graph/model"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/handler/graph/scalar"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase/output"
)

func TodoFromOutput(out *output.TodoOutput, name scalar.ResourceName) *model.Todo {
	if out == nil || out.Todo == nil {
		return nil
	}
	return &model.Todo{
		Name:      name,
		Title:     out.Todo.Title,
		Content:   out.Todo.Content,
		Status:    toGraphQLStatus(out.Todo.Status),
		DueDate:   (*scalar.Time)(out.Todo.DueDate),
		CreatedAt: scalar.Time(out.Todo.CreatedAt),
	}
}

func toGraphQLStatus(s entity.TodoStatus) model.TodoStatus {
	switch s {
	case entity.TodoStatusPENDING:
		return model.TodoStatusTodoStatusTodo
	case entity.TodoStatusIN_PROGRESS:
		return model.TodoStatusTodoStatusInProgress
	case entity.TodoStatusDONE:
		return model.TodoStatusTodoStatusDone
	default:
		return model.TodoStatusTodoStatusUnspecified
	}
}

func TodoConnectionFromOutput(out *output.TodoListOutput, listName scalar.ResourceName) *model.TodoConnection {
	if out == nil || len(out.Todos) == 0 {
		return &model.TodoConnection{
			Edges:      []*model.TodoEdge{},
			PageInfo:   &model.PageInfo{HasNextPage: false},
			TotalCount: 0,
		}
	}

	edges := make([]*model.TodoEdge, len(out.Todos))
	for i, todo := range out.Todos {
		todoName := listName
		if todo != nil && todo.ID != 0 {
			todoName = scalar.ResourceName(fmt.Sprintf("users/%d/todo-lists/%d/todos/%d", todo.CreatorID, todo.ListID, todo.ID))
		}
		edges[i] = &model.TodoEdge{
			Node: TodoFromOutput(&output.TodoOutput{Todo: todo}, todoName),
		}
	}

	return &model.TodoConnection{
		Edges:      edges,
		PageInfo:   &model.PageInfo{HasNextPage: false}, // For simplicity, we set HasNextPage to false. You can implement proper pagination logic here.
		ListName:   out.ListName,
		TotalCount: out.TotalCount,
	}
}

func ToDomainTodoStatus(status *model.TodoStatus) (*entity.TodoStatus, error) {
	if status == nil {
		return nil, nil
	}

	var out entity.TodoStatus
	switch *status {
	case model.TodoStatusTodoStatusUnspecified, model.TodoStatusTodoStatusTodo:
		out = entity.TodoStatusPENDING
	case model.TodoStatusTodoStatusInProgress:
		out = entity.TodoStatusIN_PROGRESS
	case model.TodoStatusTodoStatusDone:
		out = entity.TodoStatusDONE
	default:
		return nil, apperrors.NewInvalidParameter("invalid todo status", nil)
	}

	return &out, nil
}
