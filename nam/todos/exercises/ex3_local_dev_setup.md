# EX3: Local Dev Setup + List RPC Methods

## Part 1 — Local Development Setup

### Changes Made

Two files were created / modified to make `make run` work in Phase 1.

---

#### 1. `cmd/server/main.go` — Application Entry Point

```
gRPC request
     │
     ▼
grpc.Server  (port :50051)
     │
     └── TodosServiceServer (handler.NewServer)
              │
              ├── usecase.TodoGetter  ← service.NewTodoGetter(queriesGW)
              └── usecase.TodoUpdater ← service.NewTodoUpdater(queriesGW, commandsGW)
                                               │
                                        persistence.NewTodoQueriesGateway()
                                        persistence.NewTodoCommandsGateway()
                                        (no-op gateways; GORM added in Phase 2)
```

**What `main.go` does (4 steps):**

1. **Read config** — `SERVER_PORT` from env (default `50051`).
2. **Wire dependencies** manually: gateway constructors → services → handler.  
   *(Wire-based DI replaces this in Phase 2.)*
3. **Create the gRPC server** and register `TodosServiceServer`.
4. **Graceful shutdown** — waits for `SIGINT`/`SIGTERM`, then calls `GracefulStop()`.

---

#### 2. `internal/service/todo_updater.go` — TodoUpdater Service

Implements the **Read-Modify-Write** pattern required by the `UpdateTodo` handler:

```
input.TodoUpdater
       │
       ├─ queries.Get()          ← read current todo
       ├─ apply non-nil fields   ← partial update via pointer fields
       └─ commands.Update()      ← persist changes
              │
       output.TodoUpdater
```

Only pointer fields in `input.TodoUpdater` (`*Title`, `*Content`, `*Status`, `*DueDate`)
are applied — `nil` means "no change", mirroring the FieldMask sent by the client.

---

#### Note on `infra/persistence/todo_queries.go` and `todo_commands.go`

`NewTodoQueriesGateway` and `NewTodoCommandsGateway` are **no-op placeholders** —
they satisfy the gateway interfaces but all methods return `nil` (no real DB calls).
This is intentional for Phase 1: the goal is to verify the gRPC stack compiles and
starts correctly. Real GORM implementations are added in Phase 2.

---

### How to Run

```bash
# Start the server (runs on :50051 by default)
make run
# Output: gRPC server listening on :50051

# Use a custom port
SERVER_PORT=9090 make run
```

### Verify with grpcurl (optional)

```bash
# List all services
grpcurl -plaintext localhost:50051 list

# Describe the TodosService
grpcurl -plaintext localhost:50051 describe todo.v1.TodosService

# Call GetTodo (returns NOT_FOUND — no DB yet, Phase 2 wires GORM)
grpcurl -plaintext -d '{"name":"users/1/todo-lists/1/todos/1"}' \
  localhost:50051 todo.v1.TodosService/GetTodo
```

---

## Part 2 — TodosService RPC Methods

Source: `proto/todo/v1/todo.proto`

### Service Definition

```protobuf
service TodosService {
    rpc GetTodo    (GetTodoRequest)    returns (GetTodoResponse);
    rpc ListTodos  (ListTodosRequest)  returns (ListTodosResponse);
    rpc CreateTodo (CreateTodoRequest) returns (CreateTodoResponse);
    rpc UpdateTodo (UpdateTodoRequest) returns (UpdateTodoResponse);
    rpc DeleteTodo (DeleteTodoRequest) returns (google.protobuf.Empty);
}
```

---

### RPC Method Catalogue

#### 1. `GetTodo`

| Field | Detail |
|-------|--------|
| **Request** | `GetTodoRequest` |
| **Response** | `GetTodoResponse` |
| **Full method** | `/todo.v1.TodosService/GetTodo` |

**Request message:**
```protobuf
message GetTodoRequest {
    string name = 1;   // Resource name: "users/{user_id}/todo-lists/{list_id}/todos/{todo_id}"
}
```

**Response message:**
```protobuf
message GetTodoResponse {
    Todo todo = 1;
}
```

**Description:**  
Fetches a single Todo by its resource name.  The name encodes the full hierarchy
(`user → todo-list → todo`), following the Google AIP resource-name pattern.
Returns `NOT_FOUND` if the todo does not exist, `INVALID_ARGUMENT` if the name
format is invalid.

---

#### 2. `ListTodos`

| Field | Detail |
|-------|--------|
| **Request** | `ListTodosRequest` |
| **Response** | `ListTodosResponse` |
| **Full method** | `/todo.v1.TodosService/ListTodos` |

**Request message:**
```protobuf
message ListTodosRequest {
    string parent          = 1;  // Parent resource: "users/{user_id}/todo-lists/{list_id}"
    optional TodoStatus status = 2;  // Filter by status (omit to return all)
    int32  limit           = 3;  // Page size
    int32  offset          = 4;  // Pagination offset
}
```

**Response message:**
```protobuf
message ListTodosResponse {
    repeated Todo todos = 1;
    int32 total_count   = 2;  // Total number of matching todos (for pagination UI)
}
```

**Description:**  
Returns a paginated list of Todos under a given todo-list parent.
Supports optional status filtering and offset-based pagination.
The `total_count` field lets clients calculate the total number of pages.

---

#### 3. `CreateTodo`

| Field | Detail |
|-------|--------|
| **Request** | `CreateTodoRequest` |
| **Response** | `CreateTodoResponse` |
| **Full method** | `/todo.v1.TodosService/CreateTodo` |

**Request message:**
```protobuf
message CreateTodoRequest {
    string parent = 1;  // "users/{user_id}/todo-lists/{list_id}"
    Todo   todo   = 2;  // Todo to create (server assigns ID and timestamps)
}
```

**Response message:**
```protobuf
message CreateTodoResponse {
    Todo todo = 1;  // Created todo with server-assigned name/timestamps
}
```

**Description:**  
Creates a new Todo under the specified parent todo-list.
The server assigns the `todo.name`, `created_at`, and `updated_at` fields —
callers should leave those fields blank in the request.
title and content are required; due_date and status are optional.

---

#### 4. `UpdateTodo`

| Field | Detail |
|-------|--------|
| **Request** | `UpdateTodoRequest` |
| **Response** | `UpdateTodoResponse` |
| **Full method** | `/todo.v1.TodosService/UpdateTodo` |

**Request message:**
```protobuf
message UpdateTodoRequest {
    Todo                      todo        = 1;  // Todo with updated field values
    google.protobuf.FieldMask update_mask = 2;  // Which fields to update
}
```

**Response message:**
```protobuf
message UpdateTodoResponse {
    Todo todo = 1;  // Updated todo
}
```

**Description:**  
Performs a **partial update** using a `FieldMask`.
Only the paths listed in `update_mask` are changed; all other fields keep their
current values (Read-Modify-Write).  Supported paths: `title`, `content`,
`status`, `due_date`.  If `update_mask` is empty, all four fields are updated.

---

#### 5. `DeleteTodo`

| Field | Detail |
|-------|--------|
| **Request** | `DeleteTodoRequest` |
| **Response** | `google.protobuf.Empty` |
| **Full method** | `/todo.v1.TodosService/DeleteTodo` |

**Request message:**
```protobuf
message DeleteTodoRequest {
    string name = 1;  // Resource name: "users/{user_id}/todo-lists/{list_id}/todos/{todo_id}"
}
```

**Response message:** `google.protobuf.Empty` (no body)

**Description:**  
Deletes a Todo by its resource name and returns an empty response on success.
Returns `NOT_FOUND` if the todo does not exist.
In a production system this would be a soft-delete (GORM `DeletedAt` field),
but Phase 1 uses a hard-delete from the in-memory store.

---

### Todo Message (shared across RPCs)

```protobuf
message Todo {
    string                    name       = 1;  // "users/{uid}/todo-lists/{lid}/todos/{tid}"
    string                    title      = 2;
    string                    content    = 3;
    TodoStatus                status     = 4;
    google.protobuf.Timestamp due_date   = 5;
    google.protobuf.Timestamp created_at = 6;
    google.protobuf.Timestamp updated_at = 7;
}
```

### TodoStatus Enum

| Proto value | Domain value | Meaning |
|-------------|--------------|---------|
| `UNSPECIFIED` (0) | — | Default / not set |
| `PENDING` (1) | `TodoStatusPENDING` | Not yet started |
| `IN_PROGRESS` (2) | `TodoStatusIN_PROGRESS` | Being worked on |
| `DONE` (3) | `TodoStatusDONE` | Completed |

---

## Summary

| What | Result |
|------|--------|
| `make run` | Server starts on `:50051` — `gRPC server listening on :50051` |
| `make test` | All packages compile; no test files yet (added in Phase 4) |
| RPC methods | 5 — GetTodo, ListTodos, CreateTodo, UpdateTodo, DeleteTodo |
| Proto file | `proto/todo/v1/todo.proto` |
| Generated code | `proto/todo/v1/todo.pb.go`, `todo_grpc.pb.go` |
