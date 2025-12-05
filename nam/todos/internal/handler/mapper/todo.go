package mapper

// todo.go — Proto ↔ Domain Mapper
//
// Week 2: gRPC & Protobuf — Handler Layer
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
// See: resources/week-02-grpc-protobuf.md (mapper pattern, resource names)
