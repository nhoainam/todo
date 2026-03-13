package grpc_client

// todo_service.go — TodoService gRPC Client Implementation
//
// Phase 3: GraphQL & BFF Pattern
//
// This file implements the TodoServiceGateway interface using a gRPC client.
//
// Responsibilities:
// 1. Hold the gRPC client connection (todospb.TodosServiceClient)
// 2. Translate gateway method calls to gRPC requests
// 3. Map proto responses back to domain types
// 4. Map gRPC errors to domain AppError
//
// Pattern:
//   func (c *todoServiceClient) GetTodo(ctx context.Context, name string) (*domain.Todo, error) {
//       resp, err := c.client.GetTodo(ctx, &pb.GetTodoRequest{Name: name})
//       if err != nil {
//           return nil, mapGRPCError(err)  // gRPC status → AppError
//       }
//       return mapper.TodoFromProto(resp), nil
//   }
//
// See: resources/phase-03-graphql-bff.md (gRPC client implementation)
