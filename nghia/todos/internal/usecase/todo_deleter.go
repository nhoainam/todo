package usecase

// todo_deleter.go — DeleteTodo Use Case
//
// Week 2: gRPC & Protobuf — UseCase Layer
//
// This file is responsible for:
// 1. Define the TodoDeleter interface
// 2. Implement Execute: verify todo exists, then delete via gateway
//
// Consider: Should delete be soft-delete (mark as deleted) or hard-delete?
// In production, we often use soft-delete (GORM's DeletedAt field).
//
// See: resources/week-02-grpc-protobuf.md (use case pattern)
// See: resources/week-03-gorm-wire.md (GORM soft delete)
