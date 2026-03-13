package usecase

// todo_getter.go — BFF GetTodo Use Case
//
// Phase 3: GraphQL & BFF Pattern
//
// BFF use cases orchestrate calls to backend gRPC services.
// Unlike backend use cases that access the DB directly, BFF use cases
// call gateway interfaces that wrap gRPC clients.
//
// Pattern:
//   1. Receive input from GraphQL resolver
//   2. Call TodoServiceGateway (which calls backend gRPC)
//   3. Map response and return
//
// See: resources/phase-03-graphql-bff.md (BFF use case pattern)
