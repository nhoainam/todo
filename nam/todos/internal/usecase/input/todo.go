package input

import (
	"time"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
)

// todo.go — Use Case Input DTOs
//
// Phase 1: gRPC & Protobuf — UseCase Layer
//
// This file is responsible for:
// 1. Define input structs for each use case:
//    - TodoGetter  { TodoID domain.TodoID }
//    - TodoCreator { Title string, Content string, ListID domain.TodoListID, CreatorID domain.UserID }
//    - TodoUpdater { TodoID domain.TodoID, Title *string, Content *string, Status *domain.TodoStatus }
//    - TodoDeleter { TodoID domain.TodoID }
//    - TodoLister  { ListID domain.TodoListID, Status *domain.TodoStatus, Offset int, Limit int }
//
// Why input DTOs?
// - Decouple use case inputs from transport (proto) messages
// - The handler converts proto → input DTO
// - Use pointer fields for optional/partial updates (e.g., *string for Title in update)
//
// See: resources/phase-01-architecture-grpc.md (DTO pattern)
// See: resources/phase-01-architecture-grpc.md (input/output DTOs)

type TodoGetter struct {
	Name entity.TodoResourceName
}

type TodoCreator struct {
	Title     string
	Content   string
	DueDate   *time.Time
	ListID    entity.TodoListID
	CreatorID entity.UserID
	Status    entity.TodoStatus
}

type TodoUpdater struct {
	Name    entity.TodoResourceName
	Title   *string
	Content *string
	Status  *entity.TodoStatus
	DueDate *time.Time
}

type TodoDeleter struct {
	Name entity.TodoResourceName
}

type TodoLister struct {
	Name   entity.TodoListResourceName
	Status *entity.TodoStatus
	Limit  int32
	Offset int32
}
