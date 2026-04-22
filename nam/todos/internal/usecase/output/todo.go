package output

import "github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"

// todo.go — Use Case Output DTOs
//
// Phase 1: gRPC & Protobuf — UseCase Layer
//
// This file is responsible for:
// 1. Define output structs for each use case:
//    - TodoGetter  { Todo *domain.Todo }
//    - TodoCreator { Todo *domain.Todo }
//    - TodoUpdater { Todo *domain.Todo }
//    - TodoLister  { Todos []*domain.Todo, TotalCount int, HasNextPage bool }
//
// Output DTOs can wrap domain entities directly since the domain
// is the "innermost" layer. The handler then maps output → proto response.
//
// See: resources/phase-01-architecture-grpc.md (DTO pattern)

type TodoGetter struct {
	Todo *entity.Todo
}

type TodoCreator struct {
	Todo *entity.Todo
}

type TodoUpdater struct {
	Todo *entity.Todo
}

type TodoLister struct {
	Todos      []*entity.Todo
	ListName   string
	TotalCount int32
}
