package mapper

import (
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/infra/persistence/model"
)

var todoStatusFromInt = map[int]entity.TodoStatus{
	0: entity.TodoStatusPENDING,
	1: entity.TodoStatusIN_PROGRESS,
	2: entity.TodoStatusDONE,
}

var todoStatusToInt = map[entity.TodoStatus]int{
	entity.TodoStatusPENDING:     0,
	entity.TodoStatusIN_PROGRESS: 1,
	entity.TodoStatusDONE:        2,
}

// ToDomainTodo converts a GORM model to a domain entity.
func ToDomainTodo(m *model.Todo) *entity.Todo {
	return &entity.Todo{
		ID:        entity.TodoID(m.ID),
		ListID:    entity.TodoListID(m.ListID),
		CreatorID: entity.UserID(m.CreatorID),
		Title:     m.Title,
		Content:   m.Content,
		Status:    todoStatusFromInt[m.Status],
		Priority:  entity.Priority(m.Priority),
		DueDate:   m.DueDate,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// ToModelTodo converts a domain entity to a GORM model.
func ToModelTodo(e *entity.Todo) *model.Todo {
	return &model.Todo{
		ID:        int64(e.ID),
		ListID:    int64(e.ListID),
		CreatorID: int64(e.CreatorID),
		Title:     e.Title,
		Content:   e.Content,
		Status:    todoStatusToInt[e.Status],
		Priority:  e.Priority.Int(),
		DueDate:   e.DueDate,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
