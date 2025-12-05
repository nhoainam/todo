package domain

// todo_list.go — TodoList Entity & Strong Types
//
// Phase 1: Clean Architecture — Domain Layer
//
// This file is responsible for:
// 1. Define TodoListID as a strong type
// 2. Define the TodoList entity struct with fields:
//    - ID        TodoListID
//    - Name      string
//    - OwnerID   UserID
//    - CreatedAt time.Time
//    - UpdatedAt time.Time
//
// A TodoList groups multiple Todos together.
// Resource name format: users/{user_id}/todo-lists/{list_id}
//
// See: resources/phase-01-architecture-grpc.md (entity design)
// See: resources/phase-01-architecture-grpc.md (resource names)
