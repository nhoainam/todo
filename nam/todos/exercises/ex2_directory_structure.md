# EX2: Directory Structure Map & Clean Architecture Layers

## Overview

The `todos/` project follows **Clean Architecture** — a layered design where dependencies only point inward (toward the domain). The innermost layer (Domain) has zero external dependencies; each outer layer depends only on layers closer to the center.

```
┌────────────────────────────────────────────────────────────┐
│  Entry Point / Config  (cmd/, di/, database/, Makefile)    │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Infrastructure Layer  (internal/infra/)             │  │
│  │  Handler Layer         (internal/handler/, proto/)   │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │  Service Layer  (internal/service/)            │  │  │
│  │  │  ┌──────────────────────────────────────────┐  │  │  │
│  │  │  │  UseCase Layer  (internal/usecase/)       │  │  │  │
│  │  │  │  ┌────────────────────────────────────┐  │  │  │  │
│  │  │  │  │  Domain Layer                      │  │  │  │  │
│  │  │  │  │  (internal/domain/, apperrors/)     │  │  │  │  │
│  │  │  │  └────────────────────────────────────┘  │  │  │  │
│  │  │  └──────────────────────────────────────────┘  │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────┘
```

---

## Annotated Directory Tree

```
todos/
│
├── go.mod                          [ROOT] Go module declaration & dependency list
├── go.sum                          [ROOT] Cryptographic checksums for all dependencies
├── Makefile                        [ROOT] Developer commands: make run, make generate, etc.
├── .gitkeep                        [ROOT] Placeholder to track empty dirs in git
│
├── cmd/
│   └── server/
│       └── main.go                 [ENTRY POINT] Application bootstrap: load config,
│                                                  initialize Wire DI, start gRPC server
│
├── database/
│   └── schema.rb                   [INFRA/CONFIG] Ridgepole schema file — declarative DB
│                                                  table definitions (Phase 2)
│
├── di/
│   └── wire.go                     [ENTRY POINT] Google Wire dependency-injection wiring:
│                                                  describes how to build the full object
│                                                  graph from providers (Phase 2)
│
├── exercises/
│   ├── ex1_sequence_diagram.md     [DOCS] EX1 deliverable
│   └── ex2_directory_structure.md  [DOCS] This file — EX2 deliverable
│
├── proto/
│   └── todo/
│       └── v1/
│           ├── todo.proto          [CONTRACT] Source of truth gRPC service definition:
│           │                                  service TodosService, all RPC methods,
│           │                                  all message types, TodoStatus enum
│           ├── todo.pb.go          [GENERATED] Go structs for all proto messages
│           │                                  (auto-generated, do not edit)
│           ├── todo_grpc.pb.go     [GENERATED] Go interfaces/stubs for gRPC server &
│           │                                  client (auto-generated, do not edit)
│           └── generate.go        [BUILD] //go:generate directive that runs protoc
│
└── internal/
    │
    ├── apperrors/
    │   └── errors.go               [DOMAIN] Custom error types with error codes
    │                                        (NotFound, InvalidParameter, AuthZ, AuthN,
    │                                        Internal). Kept close to domain because
    │                                        business rules dictate which errors exist.
    │                                        Zero external dependencies.
    │
    ├── domain/
    │   │
    │   ├── entity/
    │   │   ├── todo.go             [DOMAIN] Todo entity struct + TodoID strong type +
    │   │   │                                TodoStatus enum + IsOverdue() business method.
    │   │   │                                No DB tags, no proto tags — pure domain.
    │   │   ├── todo_list.go        [DOMAIN] TodoList entity + TodoListID strong type.
    │   │   │                                Groups multiple Todos together.
    │   │   ├── user.go             [DOMAIN] UserID strong type. Todos service doesn't
    │   │   │                                own user data but references users safely.
    │   │   └── resource_name.go    [DOMAIN] TodoResourceName struct + String() +
    │   │                                    ParseTodoResourceName(). Encodes the full
    │   │                                    hierarchy: users/{id}/todo-lists/{id}/todos/{id}
    │   │
    │   └── gateway/
    │       ├── todo_commands.go    [DOMAIN] TodoCommandsGateway interface — write ops.
    │       │                                See gateway catalog below.
    │       └── todo_queries.go     [DOMAIN] TodoQueriesGateway interface — read ops.
    │                                        See gateway catalog below.
    │
    ├── usecase/
    │   ├── todo_getter.go          [USECASE] TodoGetter interface — defines the contract
    │   │                                     for fetching a single todo by resource name.
    │   ├── todo_creator.go         [USECASE] TodoCreator interface — contract for creating
    │   │                                     a new todo in a given list.
    │   ├── todo_updater.go         [USECASE] TodoUpdater interface — contract for updating
    │   │                                     todo fields (title, content, status).
    │   ├── todo_deleter.go         [USECASE] TodoDeleter interface — contract for deleting
    │   │                                     a todo by resource name.
    │   ├── todo_lister.go          [USECASE] TodoLister interface — contract for listing
    │   │                                     todos in a list with pagination/filter.
    │   ├── input/
    │   │   └── todo.go             [USECASE] Input DTOs for all use cases (e.g.,
    │   │                                     input.TodoGetter{Name: TodoResourceName}).
    │   │                                     Decouples use case contracts from proto types.
    │   └── output/
    │       └── todo.go             [USECASE] Output DTOs for all use cases (e.g.,
    │                                         output.TodoGetter{Todo: *entity.Todo}).
    │                                         The handler maps these to proto responses.
    │
    ├── service/
    │   └── todo_getter.go          [SERVICE] Concrete implementation of usecase.TodoGetter.
    │                                         Contains real business logic: calls
    │                                         TodoQueriesGateway, checks for nil, wraps
    │                                         AppError. Depends on gateway interfaces, not
    │                                         GORM or any DB driver directly.
    │
    ├── handler/
    │   ├── todo_handler.go         [HANDLER] gRPC service implementation. Implements
    │   │                                     TodosServiceServer (generated interface).
    │   │                                     Follows 5-step pattern per RPC method:
    │   │                                     Parse → Build Input → Validate → Execute
    │   │                                     → Map Response. Converts AppError to gRPC
    │   │                                     status codes. No business logic here.
    │   └── mapper/
    │       └── todo.go             [HANDLER] Proto ↔ Domain mappers. Converts proto
    │                                         messages to domain types and vice versa
    │                                         (e.g., statusToPb, TodoToProto). Tested
    │                                         independently from the handler.
    │
    └── infra/
        └── persistence/
            ├── db.go               [INFRA] GORM database connection setup +
            │                               db-from-context pattern (Phase 2).
            ├── todo_queries.go     [INFRA] Implements domain/gateway.TodoQueriesGateway
            │                               using GORM. Translates domain calls to SQL.
            ├── todo_commands.go    [INFRA] Implements domain/gateway.TodoCommandsGateway
            │                               using GORM. Handles INSERT/UPDATE/DELETE.
            ├── model/
            │   └── todo.go         [INFRA] GORM model struct with db column tags.
            │                               Separate from entity.Todo — DB schema changes
            │                               don't leak into the domain.
            └── mapper/
                └── todo.go         [INFRA] Domain entity ↔ GORM model mappers. Used
                                            by persistence implementations when reading
                                            from / writing to the database (Phase 2).
```

---

## Gateway Interface Catalog

Gateway interfaces are defined in `internal/domain/gateway/`. They are the boundary between business logic and infrastructure — the domain declares **what** it needs, and the infrastructure layer provides **how** it is done.

### 1. `TodoCommandsGateway` — `gateway/todo_commands.go`

```go
type TodoCommandsGateway interface {
    Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
    Delete(ctx context.Context, todoID entity.TodoID) error
    Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
}
```

| Method | Purpose |
|--------|---------|
| `Create` | Persists a new `entity.Todo` to storage and returns the saved record (with generated ID, timestamps filled in). |
| `Delete` | Removes the todo identified by `TodoID` from storage. Returns an error if not found or DB failure. |
| `Update` | Persists changes to an existing `entity.Todo` (title, content, status, due date) and returns the updated record. |

**Why it's in the domain layer:**
The domain defines *what write operations it needs* (`Create`, `Update`, `Delete`) without caring whether the backing store is PostgreSQL, MySQL, an in-memory map, or a mock. Any infrastructure that satisfies this interface can be injected.

**Command/Query Separation:** This interface only contains mutating operations. Read operations belong to `TodoQueriesGateway`, preventing read implementations from accidentally being used for writes and vice versa.

---

### 2. `TodoQueriesGateway` — `gateway/todo_queries.go`

```go
type GetTodoOptions struct {
    // (reserved for future options: preloads, field masks, etc.)
}

type ListTodosOptions struct {
    // (reserved for future options: pagination, filters, etc.)
}

type TodoQueriesGateway interface {
    Get(ctx context.Context, todoID entity.TodoID, opts *GetTodoOptions) (*entity.Todo, error)
    // List(ctx context.Context, opts *ListTodosOptions) (*query.OffsetPageResult[*entity.Todo], error)
}
```

| Method | Purpose |
|--------|---------|
| `Get` | Fetches a single `entity.Todo` by its `TodoID`. Returns `nil, nil` when the todo doesn't exist (the service layer then converts this to a `NotFound` AppError). The `opts` parameter allows future extension (field masks, eager loading) without breaking the interface. |
| `List` *(stub)* | (Not yet implemented) Will fetch a paginated list of todos filtered and sorted by `ListTodosOptions`. Returns an `OffsetPageResult` containing items, total count, and pagination metadata. |

**Why it's in the domain layer:**
Business rules need to query entities (e.g., "fetch this todo so we can check its owner before updating"). The domain declares the shape of the query via the interface, while the concrete SQL/GORM query lives only in `infra/persistence/`.

**Options structs:** `GetTodoOptions` and `ListTodosOptions` are defined alongside the interface (not in a separate package) so the domain layer fully owns the query contract, including extensible options.

---

## Dependency Flow Summary

```
proto/todo/v1/todo.proto
        │  (defines gRPC contract)
        ▼
internal/handler/todo_handler.go        ← implements todov1.TodosServiceServer
        │  depends on usecase interfaces
        ▼
internal/usecase/todo_getter.go         ← interface: TodoGetter
        │  implemented by
        ▼
internal/service/todo_getter.go         ← concrete impl; depends on gateway interface
        │  depends on
        ▼
internal/domain/gateway/todo_queries.go ← interface: TodoQueriesGateway
        │  implemented by
        ▼
internal/infra/persistence/todo_queries.go ← GORM query; reads/writes DB
```

The key principle: **arrows always point inward**. The domain (`entity/`, `gateway/`) never imports from `handler/`, `service/`, or `infra/`. This means the entire business logic can be tested without a database by injecting mock gateway implementations.
