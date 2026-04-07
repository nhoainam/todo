package entity

// user.go — User Strong Type
//
// Phase 1: Clean Architecture — Domain Layer
//
// This file is responsible for:
// 1. Define UserID as a strong type
//
// The Todos service doesn't own user data, but it references users
// as creators and owners. We define UserID here so the domain layer
// can reference users in a type-safe way.
//
// Example:
//   type UserID string
//   func (id UserID) String() string { return string(id) }
//
// See: resources/phase-01-architecture-grpc.md (strong typing)

type UserID int64

func (id UserID) Int64() int64 { return int64(id) }
