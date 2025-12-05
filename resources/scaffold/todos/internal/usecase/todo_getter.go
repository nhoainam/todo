package usecase

// todo_getter.go — GetTodo Use Case
//
// Phase 1: gRPC & Protobuf — UseCase Layer
//
// This file is responsible for:
// 1. Define the TodoGetter interface with an Execute method
// 2. Define todoGetterImpl struct that holds gateway dependencies:
//    - TodoQueriesGateway (for reading from DB)
// 3. Define NewTodoGetter constructor (used by Wire for DI)
// 4. Implement Execute method with business logic:
//    - Receive input DTO (input.TodoGetter)
//    - Call gateway to fetch the todo
//    - Return output DTO (output.TodoGetter) or AppError
//
// Naming convention: use cases are named as {Action}{Entity}
//   - TodoGetter (not GetTodoUseCase)
//   - TodoCreator (not CreateTodoUseCase)
//
// The use case layer:
// - Contains business rules (e.g., permission checks, validation)
// - Depends on gateway INTERFACES (not implementations)
// - Is independent of transport (no gRPC or HTTP knowledge)
//
// See: resources/phase-01-architecture-grpc.md (use case pattern)
// See: resources/phase-01-architecture-grpc.md (dependency inversion)
