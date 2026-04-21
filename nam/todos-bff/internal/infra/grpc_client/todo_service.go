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
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	return mapper.TodoFromProto(resp.Todo), nil // Map proto response to domain entity
}

func (c *todoServiceClient) ListTodos(ctx context.Context, parent string, opts *gateway.ListTodosOptions) (*gateway.OffsetPageResult[*entity.Todo], error) {
	limit := int32(20)
	offset := int32(0)
	if opts != nil && opts.Pagination != nil {
		limit = int32(opts.Pagination.Limit)
		offset = int32(opts.Pagination.Offset)
	}

	resp, err := c.client.ListTodos(ctx, &todopb.ListTodosRequest{Parent: parent, Limit: limit, Offset: offset})
	if err != nil {
		return nil, mapGRPCError(err)
	}
	return &gateway.OffsetPageResult[*entity.Todo]{
		Items:      mapper.ListTodosFromProto(resp),
		TotalCount: int64(resp.GetTotalCount()),
		ListName:   resp.GetListName(),
		Page: &gateway.OffsetPage{
			Offset: int(offset),
			Limit:  int(limit),
		},
	}, nil
}

func (c *todoServiceClient) UpdateTodo(ctx context.Context, input *gateway.UpdateTodoInput) (*entity.Todo, error) {
	if input == nil {
		return nil, apperrors.NewInvalidParameter("input cannot be nil", nil)
	}

	if input.Name == "" {
		return nil, apperrors.NewInvalidParameter("name is required", nil)
	}
	if input.Title == nil && input.Content == nil && input.Status == nil && input.DueDate == nil {
		return nil, apperrors.NewInvalidParameter("at least one field to update must be provided", nil)
	}
	todo := &todopb.Todo{Name: input.Name}
	var updateMask []string
	if input.Title != nil {
		todo.Title = *input.Title
		updateMask = append(updateMask, "title")
	}
	if input.Content != nil {
		todo.Content = *input.Content
		updateMask = append(updateMask, "content")
	}
	if input.Status != nil {
		todo.Status = mapper.TodoStatusToProto(*input.Status)
		updateMask = append(updateMask, "status")
	}
	if input.DueDate != nil {
		todo.DueDate = timestamppb.New(*input.DueDate)
		updateMask = append(updateMask, "due_date")
	}
	req := &todopb.UpdateTodoRequest{
		Todo:       todo,
		UpdateMask: &fieldmaskpb.FieldMask{Paths: updateMask},
	}
	resp, err := c.client.UpdateTodo(ctx, req)
	if err != nil {
		return nil, mapGRPCError(err)
	}
	return mapper.TodoFromProto(resp.Todo), nil
}
