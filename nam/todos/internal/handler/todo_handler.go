package handler

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/apperrors"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/handler/mapper"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/input"
	todov1 "github.com/tuannguyenandpadcojp/fresher26/nam/todos/proto/todo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// todo_handler.go — gRPC Handler for Todos Service
//
// Phase 1: gRPC & Protobuf — Handler Layer
//
// This file is responsible for:
// 1. Define a struct that implements the gRPC TodosServiceServer interface
//    (generated from proto/todo/v1/todo.proto)
// 2. Inject use case dependencies via constructor:
//    - TodoGetter, TodoCreator, TodoUpdater, TodoDeleter, TodoLister
// 3. Implement each RPC method following the 5-step handler pattern:
//
//    func (h *todosHandler) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.Todo, error) {
//        // Step 1: Parse — extract fields from the gRPC request
//        // Step 2: Build Input — create the use case input DTO
//        // Step 3: Validate — check the input (return InvalidParameter if bad)
//        // Step 4: Execute — call the use case
//        // Step 5: Map Response — convert domain entity to proto response
//    }
//
// Key principles:
// - The handler ONLY does request/response translation — no business logic
// - Use the mapper package (handler/mapper/) to convert between proto ↔ domain
// - Map AppError to gRPC status codes (NotFound → codes.NotFound, etc.)
// - Parse resource names: "users/{user_id}/todo-lists/{list_id}/todos/{todo_id}"
//
// See: resources/phase-01-architecture-grpc.md (handler 5-step pattern, error mapping)

type server struct {
	todov1.UnimplementedTodosServiceServer
	todoGetter  usecase.TodoGetter
	todoUpdater usecase.TodoUpdater
	todoLister  usecase.TodoLister
	todoCreator usecase.TodoCreator
	todoDeleter usecase.TodoDeleter
	validator   *validator.Validate
}

func NewServer(
	todoGetter usecase.TodoGetter,
	todoUpdater usecase.TodoUpdater,
	todoLister usecase.TodoLister,
	todoCreator usecase.TodoCreator,
	todoDeleter usecase.TodoDeleter,
	validator *validator.Validate,
) todov1.TodosServiceServer {
	return &server{
		todoGetter:  todoGetter,
		todoUpdater: todoUpdater,
		todoLister:  todoLister,
		todoCreator: todoCreator,
		todoDeleter: todoDeleter,
		validator:   validator,
	}
}

// toGRPCError converts an AppError to a proper gRPC status error.
// Without this, gRPC returns codes.Unknown for all domain errors.
func toGRPCError(err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case apperrors.ErrorCodeNotFound:
			return status.Error(codes.NotFound, appErr.Message)
		case apperrors.ErrorCodeInvalidParameter:
			return status.Error(codes.InvalidArgument, appErr.Message)
		case apperrors.ErrorCodeAuthZ:
			return status.Error(codes.PermissionDenied, appErr.Message)
		case apperrors.ErrorCodeAuthN:
			return status.Error(codes.Unauthenticated, appErr.Message)
		default:
			return status.Error(codes.Internal, appErr.Message)
		}
	}
	return status.Error(codes.Internal, err.Error())
}

func (s *server) GetTodo(ctx context.Context, req *todov1.GetTodoRequest) (*todov1.GetTodoResponse, error) {
	// Step 1: Parse proto request -> domain types
	name, err := entity.ParseTodoResourceName(req.Name)
	if err != nil {
		return nil, toGRPCError(apperrors.NewInvalidParameter("invalid todo resource name", err))
	}

	// Step 2: Build usecase input DTO
	in := input.TodoGetter{
		Name: *name,
	}

	// Step 3: Validate input
	if err := s.validator.Struct(&in); err != nil {
		return nil, toGRPCError(apperrors.NewInvalidParameter("invalid request", err))
	}

	// Step 4: Call usecase
	out, err := s.todoGetter.Get(ctx, &in)
	if err != nil {
		return nil, toGRPCError(err)
	}

	// Step 5: Map domain -> proto response
	return &todov1.GetTodoResponse{
		Todo: mapper.TodoToPb(out.Todo),
	}, nil
}

func (s *server) UpdateTodo(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.UpdateTodoResponse, error) {
	if req.Todo == nil {
		return nil, toGRPCError(apperrors.NewInvalidParameter("todo is required", nil))
	}

	// Step 1: Parse resource name from the todo message
	name, err := entity.ParseTodoResourceName(req.Todo.Name)
	if err != nil {
		return nil, toGRPCError(apperrors.NewInvalidParameter("invalid todo resource name", err))
	}

	// Step 2: Build input DTO — walk the FieldMask and set only the requested pointer fields
	in := input.TodoUpdater{
		Name: *name,
	}

	paths := req.GetUpdateMask().GetPaths()
	// If the mask is empty, treat it as a full update (all settable fields)
	if len(paths) == 0 {
		paths = []string{"title", "content", "status", "due_date"}
	}

	for _, path := range paths {
		switch path {
		case "title":
			v := req.Todo.Title
			in.Title = &v
		case "content":
			v := req.Todo.Content
			in.Content = &v
		case "status":
			s := mapper.PbToStatus(req.Todo.Status)
			in.Status = &s
		case "due_date":
			if req.Todo.DueDate != nil {
				t := req.Todo.DueDate.AsTime()
				in.DueDate = &t
			}
		}
	}

	// Step 3: Validate
	if err := s.validator.Struct(&in); err != nil {
		return nil, toGRPCError(apperrors.NewInvalidParameter("invalid request", err))
	}

	// Step 4: Execute
	out, err := s.todoUpdater.Update(ctx, &in)
	if err != nil {
		return nil, toGRPCError(err)
	}

	// Step 5: Map response
	return &todov1.UpdateTodoResponse{
		Todo: mapper.TodoToPb(out.Todo),
	}, nil
}
