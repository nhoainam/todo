package usecase

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/output"
)

// todo_lister.go — ListTodos Use Case
//
// Phase 1: gRPC & Protobuf — UseCase Layer
//
// This file is responsible for:
// 1. Define the TodoLister interface
// 2. Implement Execute: list todos with filtering and pagination
//
// Input should support:
//   - Filter by list ID, status, creator
//   - Pagination (offset + limit)
//
// Output should include:
//   - List of todos
//   - Pagination metadata (total count, has next page)
//
// See: resources/phase-01-architecture-grpc.md (use case pattern)
// See: resources/phase-02-database-di.md (filter and pagination patterns)

type TodoLister interface {
	List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error)
}
