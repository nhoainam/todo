package mapper

import (
	"fmt"
	"time"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
	todov1 "github.com/tuannguyenandpadcojp/fresher26/nam/todos/proto/todo/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// todo.go — Proto ↔ Domain Mapper
//
// Phase 1: gRPC & Protobuf — Handler Layer
//
// This file is responsible for:
// 1. Convert proto request messages to domain types:
//    - Proto TodoStatus enum → domain.TodoStatus
//    - Proto resource name string → domain.TodoID, domain.TodoListID, etc.
// 2. Convert domain entities to proto response messages:
//    - domain.Todo → *pb.Todo (proto message)
//    - domain.TodoList → *pb.TodoList (proto message)
// 3. Convert domain errors to gRPC status errors:
//    - domain.AppError → status.Error(codes.NotFound, msg)
//
// Why a separate mapper file?
// - Keeps the handler clean (handler calls mapper functions, not inline conversion)
// - Mappers are easily testable in isolation
// - Each layer has its own representation — mappers bridge them
//
// Example:
//   func TodoToProto(t *domain.Todo) *pb.Todo {
//       return &pb.Todo{
//           Name:   fmt.Sprintf("users/%s/todo-lists/%s/todos/%s", t.CreatorID, t.ListID, t.ID),
//           Title:  t.Title,
//           Status: statusToProto(t.Status),
//       }
//   }
//
// See: resources/phase-01-architecture-grpc.md (mapper pattern, resource names)

func statusToPb(status entity.TodoStatus) todov1.TodoStatus {
	switch status {
	case entity.TodoStatusPENDING:
		return todov1.TodoStatus_PENDING
	case entity.TodoStatusDONE:
		return todov1.TodoStatus_DONE
	case entity.TodoStatusIN_PROGRESS:
		return todov1.TodoStatus_IN_PROGRESS
	default:
		return todov1.TodoStatus_UNSPECIFIED
	}
}

func ToProtoTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func TodoToPb(t *entity.Todo) *todov1.Todo {
	if t == nil {
		return nil
	}
	return &todov1.Todo{
		Name:      fmt.Sprintf("users/%d/todo-lists/%d/todos/%d", t.CreatorID, t.ListID, t.ID),
		Title:     t.Title,
		Content:   t.Content,
		Status:    statusToPb(t.Status),
		DueDate:   ToProtoTime(t.DueDate),
		CreatedAt: ToProtoTime(&t.CreatedAt),
		UpdatedAt: ToProtoTime(&t.UpdatedAt),
	}
}
