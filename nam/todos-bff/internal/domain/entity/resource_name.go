package entity

import (
	"fmt"
	"strconv"
	"strings"
)

type TodoResourceName struct {
	UserID     UserID
	TodoListID TodoListID
	TodoID     TodoID
}

func (n TodoResourceName) String() string {
	return fmt.Sprintf("users/%d/todo-lists/%d/todos/%d", n.UserID, n.TodoListID, n.TodoID)
}

func ParseTodoResourceName(name string) (*TodoResourceName, error) {
	// Parse "users/100/todo-lists/200/todos/456" -> TodoResourceName{...}
	parts := strings.Split(name, "/")
	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid todo resource name format")
	}

	if parts[0] != "users" || parts[2] != "todo-lists" || parts[4] != "todos" {
		return nil, fmt.Errorf("invalid todo resource name: %q", name)
	}

	userIDStr := parts[1]
	todoListIDStr := parts[3]
	todoIDStr := parts[5]

	if userIDStr == "" || todoListIDStr == "" || todoIDStr == "" {
		return nil, fmt.Errorf("invalid todo resource name: %q, ids must not be empty", name)
	}

	uid, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	listID, err := strconv.ParseInt(todoListIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid todo list id: %w", err)
	}

	tid, err := strconv.ParseInt(todoIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid todo id: %w", err)
	}

	return &TodoResourceName{
		UserID:     UserID(uid),
		TodoListID: TodoListID(listID),
		TodoID:     TodoID(tid),
	}, nil
}
