package mapper

import (
	"time"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/entity"
	pb "github.com/tuannguyenandpadcojp/fresher26/nam/todos/proto/todo/v1"
)

func TodoFromProto(pbTodo *pb.Todo) *entity.Todo {
	if pbTodo == nil {
		return nil
	}

	var dueDate *time.Time
	if pbTodo.DueDate != nil {
		t := pbTodo.DueDate.AsTime()
		dueDate = &t
	}

	var createdAt time.Time
	if pbTodo.CreatedAt != nil {
		createdAt = pbTodo.CreatedAt.AsTime()
	}

	var updatedAt time.Time
	if pbTodo.UpdatedAt != nil {
		updatedAt = pbTodo.UpdatedAt.AsTime()
	}

	TodoResouceName, err := entity.ParseTodoResourceName(pbTodo.Name)
	if err != nil {
		return nil
	}

	return &entity.Todo{
		ID:        TodoResouceName.TodoID,
		ListID:    TodoResouceName.TodoListID,
		CreatorID: TodoResouceName.UserID,
		Title:     pbTodo.Title,
		Content:   pbTodo.Content,
		Status:    statusFromProto(pbTodo.Status),
		DueDate:   dueDate,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func statusFromProto(s pb.TodoStatus) entity.TodoStatus {
	switch s {
	case pb.TodoStatus_PENDING:
		return entity.TodoStatusPENDING
	case pb.TodoStatus_IN_PROGRESS:
		return entity.TodoStatusIN_PROGRESS
	case pb.TodoStatus_DONE:
		return entity.TodoStatusDONE
	default:
		return entity.TodoStatusPENDING
	}
}

func ListTodosFromProto(pbResp *pb.ListTodosResponse) []*entity.Todo {
	if pbResp == nil || len(pbResp.Todos) == 0 {
		return nil
	}

	todos := make([]*entity.Todo, len(pbResp.Todos))
	for i, pbTodo := range pbResp.Todos {
		todos[i] = TodoFromProto(pbTodo)
	}
	return todos
}

func TodoStatusToProto(s entity.TodoStatus) pb.TodoStatus {
	switch s {
	case entity.TodoStatusPENDING:
		return pb.TodoStatus_PENDING
	case entity.TodoStatusIN_PROGRESS:
		return pb.TodoStatus_IN_PROGRESS
	case entity.TodoStatusDONE:
		return pb.TodoStatus_DONE
	default:
		return pb.TodoStatus_PENDING
	}
}
