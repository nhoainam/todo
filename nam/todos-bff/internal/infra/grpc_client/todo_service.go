package grpc_client

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/apperrors"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/entity"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/gateway"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/infra/grpc_client/mapper"
	todopb "github.com/tuannguyenandpadcojp/fresher26/nam/todos/proto/todo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
//   func (c *todoServiceClient) GetTodo(ctx context.Context, name string) (*entity.Todo, error) {
//       resp, err := c.client.GetTodo(ctx, &pb.GetTodoRequest{Name: name})
//       if err != nil {
//           return nil, mapGRPCError(err)  // gRPC status → AppError
//       }
//       return mapper.TodoFromProto(resp), nil
//   }
//
// See: resources/phase-03-graphql-bff.md (gRPC client implementation)

type todoServiceClient struct {
	client todopb.TodosServiceClient
}

func NewTodoServiceClient(client todopb.TodosServiceClient) gateway.TodoServiceGateway {
	return &todoServiceClient{client: client}
}

func mapGRPCError(err error) error {
	if grpcErr, ok := status.FromError(err); ok {
		switch grpcErr.Code() {
		case codes.NotFound:
			return apperrors.NewNotFound(grpcErr.Message(), nil)
		case codes.InvalidArgument:
			return apperrors.NewInvalidParameter(grpcErr.Message(), nil)
		case codes.PermissionDenied:
			return apperrors.NewAuthZ(grpcErr.Message(), nil)
		case codes.Unauthenticated:
			return apperrors.NewAuthN(grpcErr.Message(), nil)
		default:
			return apperrors.NewInternal(grpcErr.Message(), nil)
		}
	}
	return apperrors.NewInternal(err.Error(), nil)
}

func (c *todoServiceClient) GetTodo(ctx context.Context, name string) (*entity.Todo, error) {
	resp, err := c.client.GetTodo(ctx, &todopb.GetTodoRequest{Name: name})
	if err != nil {
		return nil, mapGRPCError(err) // gRPC status → AppError
	}
	return mapper.TodoFromProto(resp), nil // Map proto response to domain entity
}
