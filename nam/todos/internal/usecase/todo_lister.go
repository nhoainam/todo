package usecase

// todo_lister.go — ListTodos Use Case
//
// Week 2: gRPC & Protobuf — UseCase Layer
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
// See: resources/week-02-grpc-protobuf.md (use case pattern)
// See: resources/week-03-gorm-wire.md (filter and pagination patterns)
