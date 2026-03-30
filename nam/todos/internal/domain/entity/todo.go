package entity

import "time"

// todo.go — Todo Entity & Strong Types
//
// Phase 1: Clean Architecture — Domain Layer
//
// This file is responsible for:
// 1. Define TodoID as a strong type (not raw string) — prevents mixing with other IDs
// 2. Define TodoStatus enum (e.g., TODO, IN_PROGRESS, DONE) with string representation
// 3. Define the Todo entity struct with fields:
//    - ID        TodoID
//    - Title     string
//    - Content   string
//    - Status    TodoStatus
//    - ListID    TodoListID  (which list this todo belongs to)
//    - CreatorID UserID      (who created it)
//    - CreatedAt time.Time
//    - UpdatedAt time.Time
// 4. Add business logic methods on the entity (e.g., IsOverdue(), MarkDone())
//
// Key principles:
// - The domain layer has ZERO external dependencies (no gRPC, no GORM, no framework imports)
// - Other layers depend on domain, never the reverse
// - Strong typing prevents bugs like passing a UserID where TodoID is expected
//
// Example strong type:
//   type TodoID string
//   func (id TodoID) String() string { return string(id) }
//
// See: resources/phase-01-architecture-grpc.md (strong typing, entity design)

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
