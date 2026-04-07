package entity

import "time"

// todo.go — BFF Domain Types
//
// Phase 3: GraphQL & BFF Pattern
//
// This file is responsible for:
// 1. Define BFF-specific domain types (may mirror backend domain)
// 2. Define strong types: TodoID, TodoListID, UserID, ResourceName
//
// The BFF has its own domain layer because:
// - It may aggregate data from multiple backend services
// - Its types may differ slightly from backend (e.g., include computed fields)
// - It maintains the same Clean Architecture principle: domain has no dependencies
//
// See: resources/phase-03-graphql-bff.md (BFF architecture)
// See: resources/phase-01-architecture-grpc.md (domain layer)

type TodoID int64

func (id TodoID) Int64() int64 { return int64(id) }

type TodoStatus string

const (
	TodoStatusPENDING     TodoStatus = "PENDING"
	TodoStatusIN_PROGRESS TodoStatus = "IN_PROGRESS"
	TodoStatusDONE        TodoStatus = "DONE"
)

type Priority int

const (
	PriorityLow    Priority = 0
	PriorityMedium Priority = 1
	PriorityHigh   Priority = 2
)

func (p Priority) Int() int { return int(p) }

type Todo struct {
	ID        TodoID
	ListID    TodoListID
	CreatorID UserID
	Title     string
	Content   string
	Status    TodoStatus
	Priority  Priority
	DueDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsOverdue checks if the todo item is past its due date and not done.
func (t *Todo) IsOverdue() bool {
	if t.DueDate == nil || t.Status == TodoStatusDONE {
		return false
	}
	return time.Now().After(*t.DueDate)
}

// MarkDone sets the status to DONE and updates the timestamp.
func (t *Todo) MarkDone() {
	t.Status = TodoStatusDONE
	t.UpdatedAt = time.Now()
}
