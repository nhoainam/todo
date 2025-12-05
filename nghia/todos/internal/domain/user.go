package domain

// user.go — User Strong Type
//
// Week 1: Clean Architecture — Domain Layer
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
// See: resources/week-01-clean-architecture.md (strong typing)
