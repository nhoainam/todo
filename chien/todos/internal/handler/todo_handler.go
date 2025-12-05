package handler

// todo_handler.go — gRPC Handler for Todos Service
//
// Week 2: gRPC & Protobuf — Handler Layer
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
// See: resources/week-02-grpc-protobuf.md (handler 5-step pattern, error mapping)
