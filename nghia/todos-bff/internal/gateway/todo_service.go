package gateway

// todo_service.go — TodoService Gateway Interface
//
// Week 4: GraphQL & BFF Pattern
//
// This file defines the interface for communicating with the backend Todos gRPC service.
//
//   type TodoServiceGateway interface {
//       GetTodo(ctx context.Context, name string) (*domain.Todo, error)
//       ListTodos(ctx context.Context, listName string, opts *ListOptions) ([]*domain.Todo, error)
//       CreateTodo(ctx context.Context, input *CreateTodoInput) (*domain.Todo, error)
//       UpdateTodo(ctx context.Context, input *UpdateTodoInput) (*domain.Todo, error)
//       DeleteTodo(ctx context.Context, name string) error
//   }
//
// The implementation (in infra/grpc_client/) translates these calls
// to actual gRPC requests to the backend service.
//
// See: resources/week-04-graphql-bff.md (gateway pattern in BFF)
