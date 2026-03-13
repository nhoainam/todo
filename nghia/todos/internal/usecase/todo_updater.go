package usecase

// todo_updater.go — UpdateTodo Use Case
//
// Phase 1: gRPC & Protobuf — UseCase Layer
//
// This file is responsible for:
// 1. Define the TodoUpdater interface
// 2. Implement Execute: fetch existing todo, apply changes, persist
//
// Pattern: Read-Modify-Write
//   1. Get current todo via TodoQueriesGateway
//   2. Apply updates from input DTO
//   3. Save via TodoCommandsGateway
//
// See: resources/phase-01-architecture-grpc.md (use case pattern)
