package domain

// todo.go — BFF Domain Types
//
// Week 4: GraphQL & BFF Pattern
//
// This file is responsible for:
// 1. Define BFF-specific domain types (may mirror backend domain)
// 2. Define strong types: TodoID, TodoListID, UserID, ResourceName
//
// The BFF has its own domain layer because:
// - It may aggregate data from multiple backend services
// - Its types may differ slightly from backend (e.g., include computed fields)
// - It maintains the same Clean Architecture principle: domain has no dependencies
//
// See: resources/week-04-graphql-bff.md (BFF architecture)
// See: resources/week-01-clean-architecture.md (domain layer)
