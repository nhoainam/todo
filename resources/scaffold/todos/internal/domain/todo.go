package domain

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
