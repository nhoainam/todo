package entity

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

	todoStatusCodeUnspecified int64 = 0
	todoStatusCodePending     int64 = 1
	todoStatusCodeInProgress  int64 = 2
	todoStatusCodeDone        int64 = 3
)

func (s TodoStatus) DBCode() (int64, error) {
	switch s {
	case TodoStatusPENDING:
		return todoStatusCodePending, nil
	case TodoStatusIN_PROGRESS:
		return todoStatusCodeInProgress, nil
	case TodoStatusDONE:
		return todoStatusCodeDone, nil
	default:
		return 0, fmt.Errorf("invalid todo status %q", s)
	}
}

func (s TodoStatus) Value() (driver.Value, error) {
	code, err := s.DBCode()
	if err != nil {
		return nil, err
	}
	return code, nil
}

func (s *TodoStatus) Scan(value interface{}) error {
	if s == nil {
		return fmt.Errorf("cannot scan TodoStatus into nil pointer")
	}

	switch v := value.(type) {
	case int64:
		mapped, err := todoStatusFromCode(v)
		if err != nil {
			return err
		}
		*s = mapped
		return nil
	case int32:
		return s.Scan(int64(v))
	case int:
		return s.Scan(int64(v))
	case []byte:
		return s.scanFromString(string(v))
	case string:
		return s.scanFromString(v)
	case nil:
		return fmt.Errorf("cannot scan NULL into TodoStatus")
	default:
		return fmt.Errorf("unsupported scan type %T for TodoStatus", value)
	}
}

func (s *TodoStatus) scanFromString(raw string) error {
	normalized := strings.TrimSpace(strings.ToUpper(raw))
	if normalized == "" {
		return fmt.Errorf("invalid empty todo status")
	}

	if code, err := strconv.ParseInt(normalized, 10, 64); err == nil {
		mapped, mapErr := todoStatusFromCode(code)
		if mapErr != nil {
			return mapErr
		}
		*s = mapped
		return nil
	}

	switch TodoStatus(normalized) {
	case TodoStatusPENDING, TodoStatusIN_PROGRESS, TodoStatusDONE:
		*s = TodoStatus(normalized)
		return nil
	default:
		return fmt.Errorf("invalid todo status %q", raw)
	}
}

func todoStatusFromCode(code int64) (TodoStatus, error) {
	switch code {
	case todoStatusCodeUnspecified, todoStatusCodePending:
		return TodoStatusPENDING, nil
	case todoStatusCodeInProgress:
		return TodoStatusIN_PROGRESS, nil
	case todoStatusCodeDone:
		return TodoStatusDONE, nil
	default:
		return "", fmt.Errorf("invalid todo status code %d", code)
	}
}

type Priority int

const (
	PriorityLow    Priority = 0
	PriorityMedium Priority = 1
	PriorityHigh   Priority = 2
)

func (p Priority) Int() int { return int(p) }

type Todo struct {
	ID        TodoID
	ListID    TodoListID `gorm:"column:todo_list_id"`
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
