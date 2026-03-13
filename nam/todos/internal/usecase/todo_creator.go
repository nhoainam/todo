package usecase

// todo_creator.go — CreateTodo Use Case
//
// Phase 1: gRPC & Protobuf — UseCase Layer
//
// This file is responsible for:
// 1. Define the TodoCreator interface
// 2. Define todoCreatorImpl struct with dependencies:
//    - TodoCommandsGateway (for writing to DB)
//    - IDGenerator (for creating new TodoID)
//    - Clock (for timestamps)
// 3. Implement Execute method:
//    - Receive input DTO (input.TodoCreator)
//    - Generate new ID
//    - Build domain entity
//    - Call gateway to persist
//    - Return output DTO
//
// Note: The creator depends on TodoCommandsGateway (write),
// while the getter depends on TodoQueriesGateway (read).
// This is the Commands/Queries separation pattern.
//
// See: resources/phase-01-architecture-grpc.md (use case pattern)
// See: resources/phase-02-database-di.md (gateway Commands/Queries separation)
