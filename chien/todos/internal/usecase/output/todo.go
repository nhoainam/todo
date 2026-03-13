package output

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
