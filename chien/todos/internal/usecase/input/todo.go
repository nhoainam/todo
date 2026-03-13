package input

// todo.go — Use Case Input DTOs
//
// Phase 1: gRPC & Protobuf — UseCase Layer
//
// This file is responsible for:
// 1. Define input structs for each use case:
//    - TodoGetter  { TodoID domain.TodoID }
//    - TodoCreator { Title string, Content string, ListID domain.TodoListID, CreatorID domain.UserID }
//    - TodoUpdater { TodoID domain.TodoID, Title *string, Content *string, Status *domain.TodoStatus }
//    - TodoDeleter { TodoID domain.TodoID }
//    - TodoLister  { ListID domain.TodoListID, Status *domain.TodoStatus, Offset int, Limit int }
//
// Why input DTOs?
// - Decouple use case inputs from transport (proto) messages
// - The handler converts proto → input DTO
// - Use pointer fields for optional/partial updates (e.g., *string for Title in update)
//
// See: resources/phase-01-architecture-grpc.md (DTO pattern)
// See: resources/phase-01-architecture-grpc.md (input/output DTOs)
