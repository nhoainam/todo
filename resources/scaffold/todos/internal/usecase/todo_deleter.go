package usecase

// todo_deleter.go — DeleteTodo Use Case
//
// Phase 1: gRPC & Protobuf — UseCase Layer
//
// This file is responsible for:
// 1. Define the TodoDeleter interface
// 2. Implement Execute: verify todo exists, then delete via gateway
//
// Consider: Should delete be soft-delete (mark as deleted) or hard-delete?
// In production, we often use soft-delete (GORM's DeletedAt field).
//
// See: resources/phase-01-architecture-grpc.md (use case pattern)
// See: resources/phase-02-database-di.md (GORM soft delete)
