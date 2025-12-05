# Phase 1: Architecture & gRPC (Days 1-3)

## Goals
- Understand the overall project architecture
- Master the Clean Architecture pattern as applied in the project
- Set up a local development environment
- Understand gRPC service definitions and code generation
- Master the gRPC interceptor chain
- Know how to implement a new gRPC endpoint

---

## 1. Interfaces in Go

### Problem Solved
In a large system, components depending directly on each other create tight coupling. Changing one component requires modifying all other components that depend on it.

**Real-world example:**
```go
// WITHOUT interfaces — tight coupling
type TodoService struct {
    reader *PostgresTodoReader  // Depends directly on Postgres
}

func (s *TodoService) GetTodo(id int64) (*Todo, error) {
    return s.reader.QueryByID(id)  // Cannot test without Postgres
}

// Problem 1: Want to write a unit test? Must set up a real Postgres instance
// Problem 2: Want to switch to MySQL? Must modify TodoService
// Problem 3: Want to mock the reader for error case testing? Not possible
```

```go
// WITH interfaces — loose coupling
type TodoService struct {
    reader gateway.TodoQueriesGateway  // Depends on interface
}

// Benefit 1: Unit tests use mocks, no DB needed
// Benefit 2: Switching databases? Just implement a new interface
// Benefit 3: Testing error cases? Mock returns error
```

### Core Concept
Go interfaces enable **dependency inversion** — high-level modules do not depend on low-level modules; both depend on abstractions.

```go
// Domain layer defines "what is needed" (WHAT)
type TodoQueriesGateway interface {
    Get(ctx context.Context, todoID entity.TodoID, opts *GetTodoOptions) (*entity.Todo, error)
    List(ctx context.Context, opts *ListTodosOptions) (*query.OffsetPageResult[*entity.Todo], error)
}

// Infrastructure layer implements "how it's done" (HOW)
type todoReader struct{}

func (r *todoReader) Get(ctx context.Context, todoID entity.TodoID, opts *gateway.GetTodoOptions) (*entity.Todo, error) {
    db, err := DBFromContext(ctx)
    // ... GORM query
}
```

### When to Use
- When separating business logic from implementation details (database, external APIs)
- When mocking dependencies for testing
- When a contract can have multiple implementations

### When NOT to Use
- For simple utility functions that don't need abstraction
- When there is only one implementation and no need for test isolation (rare in this project)

---

## 2. Clean Architecture

### Problem Solved
Without clear architecture:
- Business logic gets mixed with database queries, HTTP handlers, protobuf conversions
- Cannot test business logic independently
- Changing the database/framework affects the entire codebase
- New team members struggle to understand what the code does

**Real-world example:**
```go
// WITHOUT architecture — everything mixed in one function
func GetTodoHandler(w http.ResponseWriter, r *http.Request) {
    // HTTP parsing + DB query + business logic + response mapping — all in one place
    id := r.URL.Query().Get("id")
    db, _ := sql.Open("postgres", "host=localhost ...")
    row := db.QueryRow("SELECT * FROM todos WHERE id = $1", id)
    var todo Todo
    row.Scan(&todo.ID, &todo.Title, &todo.Status)

    // Business logic mixed with HTTP
    if todo.IsOverdue() {
        todo.Status = "overdue"
    }

    json.NewEncoder(w).Encode(todo)
}
// Problem: Want to switch from HTTP to gRPC? Must rewrite everything
// Problem: Want to test business logic? Must set up HTTP server + DB
// Problem: Want to switch from Postgres to MySQL? Must modify every handler
```

With Clean Architecture, each concern is separated and can be tested/changed independently.

### Core Concept

The project applies 4 layers with dependencies pointing inward:

```
┌─────────────────────────────────────────┐
│  Handler Layer (gRPC handlers, mappers) │ ← Knows about proto, HTTP
├─────────────────────────────────────────┤
│  UseCase Layer (interfaces + DTOs)      │ ← Defines "what needs to be done"
├─────────────────────────────────────────┤
│  Domain Layer (entities, gateways)      │ ← Center, depends on nothing
├─────────────────────────────────────────┤
│  Infrastructure Layer (GORM, gRPC)      │ ← Implements gateway interfaces
└─────────────────────────────────────────┘
```

**Dependency Rule**: Code may only depend on the layer inside it, NEVER on outer layers.

#### Layer 1: Domain Layer — `internal/domain/`

The center of the system. Defines:
- **Entity** (`domain/entity/`): Domain models with business rules

```go
// File: todos/internal/domain/entity/todo.go
type TodoID int64

type Todo struct {
    ID          TodoID
    TodoListID  *TodoListID
    Title       string
    Description *string
    Status      TodoStatus     // Pending, InProgress, Done
    Priority    Priority       // Low, Medium, High, Urgent
    DueDate     *time.Time
    AssigneeID  *UserID
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt // Soft delete
    // ...
}

// Business logic belongs in the entity
func (t *Todo) IsOverdue() bool {
    if t.DueDate == nil || t.Status == TodoStatusDone {
        return false
    }
    return time.Now().After(*t.DueDate)
}
```

- **Gateway interfaces** (`domain/gateway/`): Contracts for data access

```go
// File: todos/internal/domain/gateway/todo.go

// Separate Commands (write) and Queries (read) — CQRS-lite pattern
type TodoCommandsGateway interface {
    Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
    Delete(ctx context.Context, todoID entity.TodoID) error
    Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
}

type TodoQueriesGateway interface {
    Get(ctx context.Context, todoID entity.TodoID, opts *GetTodoOptions) (*entity.Todo, error)
    List(ctx context.Context, opts *ListTodosOptions) (*query.OffsetPageResult[*entity.Todo], error)
}
```

- **Query helpers** (`domain/query/`): Pagination, sorting abstractions

#### Layer 2: UseCase Layer — `internal/usecase/todos/`

Defines WHAT the system needs to do, NOT HOW.

```go
// File: todos/internal/usecase/todos/todo.go

// Each use case is a separate interface (Interface Segregation)
type TodoGetter interface {
    Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error)
}

type TodoDeleter interface {
    Delete(ctx context.Context, in *input.TodoDeleter) (*output.TodoDeleter, error)
}

type TodoCreator interface {
    Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error)
}
```

Input/Output DTOs are separate from entities:
```go
// File: todos/internal/usecase/todos/input/todo.go
type TodoGetter struct {
    Name todos.TodoResourceName
}

// File: todos/internal/usecase/todos/output/todo.go
type TodoGetter struct {
    Todo *todos.Todo
}
```

#### Layer 3: Service Layer (UseCase Implementation) — `internal/service/`

Implements business logic, orchestrates gateways:

```go
// File: todos/internal/service/todos/todo_getter.go

type todoGetter struct {
    todoQueriesGateway gateway.TodoQueriesGateway  // Interface, not implementation
    todoHelper         sharedport.TodoHelper
    clock              timeutil.Clock
    binder             gateway.Binder
    cfg                *config.Config
}

func NewTodoGetter(
    todoQueriesGateway gateway.TodoQueriesGateway,
    // ... other dependencies (all interfaces)
) usecase.TodoGetter {
    return &todoGetter{...}
}

func (g *todoGetter) Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error) {
    ctx = g.binder.Bind(ctx)

    // 1. Get entity from gateway
    todoEnt, err := g.todoHelper.Get(ctx, in.Name)
    if err != nil { return nil, err }

    if todoEnt == nil {
        return nil, errors.NewNotFound("todo not found", nil,
            errors.String("name", in.Name.String()),
        )
    }

    // 2. Convert entity to domain model
    todo, err := helper.TodoFromEntity(ctx, todoEnt, g.cfg, g.clock)
    if err != nil { return nil, fmt.Errorf("convert todo: %w", err) }

    return &output.TodoGetter{Todo: todo}, nil
}
```

#### Layer 4: Handler Layer — `internal/handler/`

Translates between transport protocols (gRPC/GraphQL) and use cases:

```go
// File: todos/internal/handler/grpc/service/todos/get_todo.go

func (s *service) GetTodo(
    ctx context.Context,
    req *todospb.GetTodoRequest,
) (*todospb.GetTodoResponse, error) {
    // 1. Parse proto request -> domain types
    name, err := todos.ParseTodoResourceName(req.Name)
    if err != nil {
        return nil, errors.NewInvalidParameter("invalid todo resource name", err)
    }

    // 2. Build usecase input
    in := input.TodoGetter{
        Name: cast.Value(name),
    }

    // 3. Validate
    if err := s.validator.Struct(&in); err != nil {
        return nil, errors.NewInvalidParameter("invalid request", err)
    }

    // 4. Call usecase
    out, err := s.todoGetter.Get(ctx, &in)
    if err != nil { return nil, fmt.Errorf("get todo: %w", err) }

    // 5. Map domain -> proto response
    return &todospb.GetTodoResponse{
        Todo: mapper.TodoToPb(out.Todo),
    }, nil
}
```

### When to Use Clean Architecture
- Large systems with multiple developers
- Need to test business logic independently from database/external services
- Need to change infrastructure (swap databases, change API protocols) without affecting business logic

### When NOT to Use
- Small prototypes, throwaway code
- Simple CRUD without business logic (but in this project, almost every feature has business logic)

---

## 3. Mapper Pattern

### Problem Solved
Each layer has its own data model (Proto, Entity, GraphQL Model). Sharing a single model means changes in one layer affect all other layers.

**Real-world example:**
```go
// WITHOUT mappers — using one struct for everything
type Todo struct {
    ID          int64  `json:"id" gorm:"primaryKey" protobuf:"varint,1"`
    Title       string `json:"title" gorm:"column:title" protobuf:"bytes,2"`
    Description string `json:"description" gorm:"column:description" protobuf:"bytes,3"`
    // ...40 tags for 3 different layers
}

// Problem 1: Proto needs a "is_overdue" field but DB has no such column
//            -> Must add `gorm:"-"` tag to skip, code keeps growing
// Problem 2: DB renames column from "title" to "display_name"
//            -> Must change json tag, proto tag, all API responses change
// Problem 3: Frontend doesn't need "assignee_id" (internal) but DB does
//            -> Must manually filter fields when serializing
```

The Mapper pattern allows each layer to have its own model, changed independently.

### Core Concept
Mapper functions convert between models:

```
Proto (todospb.Todo) <--mapper--> Domain (todos.Todo) <--mapper--> Entity (entity.Todo)
```

```go
// File: todos/internal/handler/grpc/mapper/todo.go

// Domain -> Proto
func TodoToPb(in *todos.Todo) *todospb.Todo {
    if in == nil { return nil }

    return &todospb.Todo{
        Name:        in.Name.String(),
        Title:       in.Title,
        Description: in.Description,
        Status:      TodoStatusToPb(in.Status),
        Priority:    PriorityToPb(in.Priority),
        DueDate:     timestamppb.New(in.DueDate),
        CreatedAt:   timestamppb.New(in.CreatedAt),
    }
}
```

### When to Use
- When converting between 2 different layers (proto <-> domain, graphql <-> domain)
- When data structures differ between layers

### When NOT to Use
- Within the same layer, between functions using the same model
- When 2 models are identical (but still consider separating to protect boundaries)

---

## 4. Strong Typing for IDs

### Problem Solved
Using `int64` for every type of ID (TodoID, TodoListID, UserID) makes it easy to mix up parameters. The compiler cannot catch passing `userID` where `todoID` is expected.

**Real-world example:**
```go
// WITHOUT strong typing — everything is int64
func GetTodo(todoListID int64, todoID int64) (*Todo, error) { ... }
func CheckPermission(userID int64, todoID int64) (bool, error) { ... }

// Calling — BUG but the compiler DOES NOT report an error!
userID := int64(100)
todoListID := int64(200)
todoID := int64(300)

GetTodo(todoID, todoListID)          // Wrong! Passing todoID where todoListID expected
CheckPermission(todoListID, userID)  // Wrong! Passing todoListID where userID expected
// Compiler sees: int64, int64 → OK ✓ (but the logic is wrong!)
```

```go
// WITH strong typing — compiler catches errors
func GetTodo(todoListID entity.TodoListID, todoID entity.TodoID) (*Todo, error) { ... }

GetTodo(entity.TodoID(300), entity.TodoListID(200))
// Compiler error: cannot use entity.TodoID as entity.TodoListID ✗
// Bug caught at compile time!
```

### Core Concept
Create custom types for each kind of ID:

```go
// File: todos/internal/domain/entity/todo.go
type TodoID int64

func (id TodoID) Int64() int64  { return int64(id) }
func (id TodoID) String() string { return strconv.FormatInt(id.Int64(), 10) }

// File: todos/internal/domain/entity/user.go
type UserID int64

// Compiler will error if you pass UserID where TodoID is expected
```

### When to Use
- For all types of IDs in the domain layer
- For enums (TodoStatus, Priority, TagType)

### When NOT to Use
- For simple values without domain meaning (string config values, etc.)

---

## 5. Configuration with envconfig

### Problem Solved
Hard-coding config values (DB host, port, API keys) directly in code prevents changes when deploying to different environments (staging, production).

**Real-world example:**
```go
// WITHOUT config management
func connectDB() *sql.DB {
    db, _ := sql.Open("postgres", "host=localhost port=5432 user=admin password=secret dbname=todos")
    return db
}
// Problem 1: Deploy to staging? Must modify code and re-compile
// Problem 2: Password is in source code → security risk
// Problem 3: Each developer has a different DB host → conflicts
```

```go
// WITH envconfig — reads from environment variables
// Local: DB_HOST=localhost DB_PORT=5432 (.env file)
// Staging: DB_HOST=staging-db.internal DB_PORT=5432 (K8s secrets)
// Production: DB_HOST=prod-db.internal DB_PORT=5432 (K8s secrets)
// Same binary, runs in different environments just by changing env vars
```

### Core Concept
Use environment variables, parsed into Go structs:

```go
// File: todos/internal/config/config.go
type Config struct {
    BaseConfig
    ServerPort            int  `envconfig:"SERVER_PORT"`
    GRPCReflectionEnabled bool `envconfig:"GRPC_REFLECTION_ENABLED"`
}

type BaseConfig struct {
    ServiceName string
    Environment Env    `envconfig:"ENVIRONMENT"`
    DB          *DBConfig
    // ...
}

func Load() (*Config, error) {
    var cfg Config
    cfg.ServiceName = "andpad-todos-api"
    err := envconfig.Process("", &cfg)  // Read from env vars
    if err != nil { return nil, fmt.Errorf("load config from env: %w", err) }
    if err := cfg.Validate(); err != nil { return nil, err }
    return &cfg, nil
}
```

### When to Use
- Any configuration that can change per environment
- Database credentials, API keys, feature flags, timeouts

### When NOT to Use
- Constants that never change (e.g., file path prefixes, enum values)

---

## 6. Protocol Buffers (Protobuf)

### Problem Solved
REST/JSON APIs:
- No formal contract between client and server
- JSON parsing is slow and not type-safe
- No built-in code generation for multiple languages
- Schema changes can easily become breaking changes if not managed carefully

**Real-world example:**
```go
// REST/JSON — no contract
// Backend returns:
{"todo": {"id": 123, "title": "Buy groceries", "due_date": "2025-03-01"}}

// Frontend developer reads API docs (if they exist), writes manually:
type Todo struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    DueDate string `json:"dueDate"`  // BUG! Backend uses "due_date", frontend uses "dueDate"
}
// This error is only discovered at runtime — no compile-time check

// Backend adds a new "priority" field — frontend doesn't know, doesn't update
// Backend renames "due_date" to "deadline" — frontend breaks immediately
```

```protobuf
// Protobuf — formal contract
message Todo {
    int64 id = 1;
    string title = 2;
    google.protobuf.Timestamp due_date = 3;
}
// Generates code for both Go and TypeScript — ALWAYS in sync
// Adding a new field? New field number (4) — old clients still work (backward compatible)
// Compiler auto-generates struct with correct field names
```

### Core Concept
Protobuf is a schema definition language that generates code for multiple languages.

Example proto definition:
```protobuf
service TodosService {
    rpc GetTodo(GetTodoRequest) returns (GetTodoResponse);
    rpc ListTodos(ListTodosRequest) returns (ListTodosResponse);
    rpc CreateTodo(CreateTodoRequest) returns (CreateTodoResponse);
    rpc DeleteTodo(DeleteTodoRequest) returns (DeleteTodoResponse);
    rpc UpdateTodoStatus(UpdateTodoStatusRequest) returns (UpdateTodoStatusResponse);
}

message GetTodoRequest {
    string name = 1;  // Resource name: "users/{id}/todo-lists/{id}/todos/{id}"
}

message GetTodoResponse {
    Todo todo = 1;
}
```

Generated Go code automatically creates:
- Struct types for each message
- Interfaces for each service (server + client)
- Marshal/Unmarshal methods

### When to Use
- Service-to-service communication (backend microservices)
- Need strong typing, code generation, performance
- Multi-language environments (Go, Java, Python, etc.)

### When NOT to Use
- Public-facing APIs for browser/mobile (use GraphQL or REST)
- Simple request/response without schema evolution needs
- Debug-friendly communication (protobuf is binary, hard to read)

---

## 7. gRPC Framework

### Problem Solved
HTTP REST:
- No standard way to define errors, pagination, streaming
- Every team does it differently
- No built-in middleware/interceptor chain

**Real-world example:**
```
// REST — every team does it differently
// Team A error:  {"error": "not found"}
// Team B error:  {"code": 404, "message": "Not Found"}
// Team C error:  {"errors": [{"field": "id", "msg": "invalid"}]}
// Frontend must handle 3 different formats!

// Pagination:
// Team A: /todos?page=1&per_page=20
// Team B: /todos?offset=0&limit=20
// Team C: /todos?cursor=abc123
```

```
// gRPC — standard for everything
// Error: always a gRPC status code (NotFound, InvalidArgument, PermissionDenied)
// Client just checks status.Code(err) — consistent across every service
// Pagination: defined in proto messages — every service follows the same pattern
```

### Core Concept

gRPC provides:
- **Strongly typed contracts** from proto definitions
- **Interceptor chain** (similar to HTTP middleware)
- **Streaming** support (unary, server-stream, client-stream, bidirectional)
- **Standard error codes** (InvalidArgument, NotFound, PermissionDenied...)

#### gRPC Server Setup

```go
// File: todos/internal/handler/grpc/grpc.go

func NewServer(
    healthService healthpb.HealthServer,
    todosService todospb.TodosServiceServer,
    adminService adminpb.AdminServiceServer,
    // ... interceptor dependencies
) (*grpc.Server, func(), error) {
    server := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            grpctrace.UnaryServerInterceptor(...),     // 1. Datadog tracing
            logging.UnaryServerInterceptor(...),        // 2. Logging
            recovery.UnaryServerInterceptor(...),       // 3. Panic recovery
            sentryinterceptor.UnaryServerInterceptor(), // 4. Sentry error reporting
            authninterceptor.UnaryServerInterceptor(), // 5. Authentication
            authzinterceptor.UnaryServerInterceptor(), // 6. Authorization
        ),
    )

    // Register services
    todospb.RegisterTodosServiceServer(server, todosService)
    adminpb.RegisterAdminServiceServer(server, adminService)
    healthpb.RegisterHealthServer(server, healthService)

    return server, cleanup, nil
}
```

### When to Use
- Internal microservice communication
- Need an interceptor chain for cross-cutting concerns (auth, logging, tracing)
- High-throughput, low-latency communication

### When NOT to Use
- Browser-facing APIs (browsers don't support gRPC natively — need gRPC-Web or a gateway)
- Large file uploads (gRPC has message size limits)

---

## 8. gRPC Interceptors

### Problem Solved
Every gRPC handler needs: authentication, authorization, logging, error handling, tracing. You can't copy-paste this code into every handler.

**Real-world example:**
```go
// WITHOUT interceptors — copy-paste auth/logging into EVERY handler
func (s *service) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.GetTodoResponse, error) {
    // Auth check — copied into every handler
    user, err := authenticateUser(ctx)
    if err != nil { return nil, err }
    if !user.HasPermission("todo.read") { return nil, status.Error(codes.PermissionDenied, "no access") }

    // Logging — copied into every handler
    log.Printf("GetTodo called by user %s", user.ID)
    defer log.Printf("GetTodo completed")

    // ... business logic
}

func (s *service) DeleteTodo(ctx context.Context, req *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
    // COPY-PASTE the same auth check — 100% identical to above
    user, err := authenticateUser(ctx)
    if err != nil { return nil, err }
    if !user.HasPermission("todo.delete") { return nil, status.Error(codes.PermissionDenied, "no access") }

    // COPY-PASTE the same logging
    log.Printf("DeleteTodo called by user %s", user.ID)
    defer log.Printf("DeleteTodo completed")

    // ... business logic
}
// 30 handlers → 30 copy-pastes of auth + logging + error handling
// Forget one handler → security hole!
```

Interceptors solve this — write once, apply to ALL handlers.

### Core Concept
Interceptors are middleware for gRPC. They run before/after the handler, in the defined order.

```
Request → Tracing → Logging → Recovery → Sentry → Auth(n) → Auth(z) → Handler → Response
                                                                          ↓
                                                                  Business Logic
```

#### Authentication Interceptor
```go
// File: todos/internal/handler/grpc/interceptor/authn/authn.go

// Extract user from gRPC metadata
func UnaryServerInterceptor(userGetter usecase.UserGetter) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
        // 1. Extract user ID from metadata header
        md, _ := metadata.FromIncomingContext(ctx)
        userID := md.Get("x-authenticated-user-id")

        // 2. Load full user object
        user, err := userGetter.Get(ctx, userID)

        // 3. Inject user into context
        ctx = auth.WithCurrentUser(ctx, user)

        // 4. Continue to next interceptor/handler
        return handler(ctx, req)
    }
}
```

#### Authorization Interceptor
Checks whether the user has permission to call this RPC method.

#### Selector Interceptor
Allows applying interceptors conditionally (e.g., skip auth for health checks).

```go
// Skip authn for health check
selector.UnaryServerInterceptor(
    authninterceptor.UnaryServerInterceptor(userGetter),
    selector.MatchFunc(func(ctx context.Context, callMeta interceptors.CallMeta) bool {
        return callMeta.FullMethod() != "/grpc.health.v1.Health/Check"
    }),
)
```

### When to Use
- Cross-cutting concerns that apply to many/all RPC methods
- Authentication, authorization, logging, tracing, error handling, rate limiting

### When NOT to Use
- Logic specific to a single RPC method — put it directly in the handler
- Logic that depends on the request body (interceptors typically only access metadata)

---

## 9. gRPC Handler Pattern

### Problem Solved
Handlers need to: parse requests, validate, call business logic, map responses. A consistent pattern is needed so everyone follows it.

**Real-world example:**
```go
// WITHOUT a pattern — every developer does it differently
// Developer A: validates first, then parses
func (s *service) GetTodo(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    if req.Name == "" { return nil, fmt.Errorf("name required") }  // validate first
    name := parseName(req.Name)                                     // parse after
    todo := s.db.Find(name.ID)                                     // calls DB directly!
    return &pb.Response{Title: todo.Title}, nil                    // inline mapping
}

// Developer B: no validation, business logic in handler
func (s *service) DeleteTodo(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    id, _ := strconv.Atoi(req.Id)                // Doesn't handle parse error!
    s.db.Delete(id)                               // Business logic in handler
    s.notifySlack("todo deleted")                 // Side effect in handler
    return &pb.Response{}, nil
}

// Problem: Code review is hard, every handler is different, easy to miss validation
```

A consistent pattern (Parse → Build Input → Validate → Execute → Map) ensures every handler is consistent and easy to review.

### Core Concept

Every handler in the project follows 5 steps:

```go
// File: todos/internal/handler/grpc/service/todos/get_todo.go

func (s *service) GetTodo(
    ctx context.Context,
    req *todospb.GetTodoRequest,
) (*todospb.GetTodoResponse, error) {
    // Step 1: Parse proto request -> domain types
    name, err := todos.ParseTodoResourceName(req.Name)
    if err != nil {
        return nil, errors.NewInvalidParameter("invalid todo resource name", err,
            errors.String("name", req.Name),
        )
    }

    // Step 2: Build usecase input DTO
    in := input.TodoGetter{
        Name: cast.Value(name),
    }

    // Step 3: Validate input
    if err := s.validator.Struct(&in); err != nil {
        return nil, errors.NewInvalidParameter("invalid request", err)
    }

    // Step 4: Call usecase
    out, err := s.todoGetter.Get(ctx, &in)
    if err != nil {
        return nil, fmt.Errorf("get todo: %w", err)
    }

    // Step 5: Map domain -> proto response
    return &todospb.GetTodoResponse{
        Todo: mapper.TodoToPb(out.Todo),
    }, nil
}
```

**Pattern**: Parse → Build Input → Validate → Execute → Map Response

### When to Use
- Every gRPC handler ALWAYS follows this pattern
- Keep handlers thin — business logic belongs in usecase/service

### When NOT to Use
- DO NOT put business logic in handlers
- DO NOT access the database directly from handlers

---

## 10. Resource Names (Google AIP Pattern)

### Problem Solved
Using raw IDs (`todoID = 123`) doesn't show where the resource belongs or what context it's in. Also hard to validate and route.

**Real-world example:**
```go
// Raw IDs — no idea where the resource belongs
GetTodo(todoID: 456)
// Which user does Todo 456 belong to? Which TodoList? No idea!
// Must query additionally: SELECT user_id, todo_list_id FROM todos WHERE id = 456
// Or pass extra parameters: GetTodo(userID: 100, todoListID: 200, todoID: 456)

// In logs:
// "Failed to get todo 456" — which user? which list?
```

```go
// Resource names — the identifier itself provides context
GetTodo(name: "users/100/todo-lists/200/todos/456")
// Immediately know: todo 456 belongs to todo-list 200 of user 100
// In logs: "Failed to get users/100/todo-lists/200/todos/456" — clear, easy to debug

// Validation: "users/abc/todo-lists/200/todos/456" → error because "abc" is not a number
// Hierarchy: know which list a todo belongs to, which user owns the list
```

### Core Concept
The project uses Google AIP-style resource names:

```
users/{user_id}/todo-lists/{list_id}/todos/{todo_id}
users/{user_id}/todo-lists/{list_id}
users/{user_id}
```

```go
// File: todos/internal/domain/model/todos/

type TodoResourceName struct {
    UserID     UserID
    TodoListID TodoListID
    TodoID     TodoID
}

func (n TodoResourceName) String() string {
    return fmt.Sprintf("users/%d/todo-lists/%d/todos/%d", n.UserID, n.TodoListID, n.TodoID)
}

func ParseTodoResourceName(name string) (*TodoResourceName, error) {
    // Parse "users/100/todo-lists/200/todos/456" -> TodoResourceName{...}
}
```

### When to Use
- All API requests use resource names instead of raw IDs
- When encoding parent-child relationships in identifiers

### When NOT to Use
- Internal database queries (use raw IDs)
- Logging (can use both — resource name for context, raw ID for lookup)

---

## 11. Error Handling Pattern

### Problem Solved
Default Go errors only have a message string. We need:
- Error codes for clients to handle each error type
- Localized messages for end users
- Structured metadata for debugging
- Mapping from domain errors to gRPC status codes

**Real-world example:**
```go
// Go standard errors — only a message
err := fmt.Errorf("todo not found")

// Client receives: "todo not found"
// Problem 1: Is this not found or permission denied? Client can't tell
// Problem 2: Which todo? What user ID? No metadata
// Problem 3: What gRPC status code? Default is Unknown — frontend can't handle
// Problem 4: How does Sentry group these? Each error message is different → 1000 issues
```

```go
// AppError — structured, with error code
err := errors.NewNotFound("todo not found", nil,
    errors.String("name", "users/100/todo-lists/200/todos/456"),
    errors.String("method", "GetTodo"),
)
// Client receives: gRPC status NotFound (404) → frontend shows "Todo not found" page
// Metadata: {name: "users/100/todo-lists/200/todos/456"} → easy to debug
// Sentry: grouped by Reason (NotFound) + method → 1 issue for all not found errors
```

### Core Concept

```go
// File: todos/internal/errors/errors.go

type AppError struct {
    Reason   Reason          // NotFound, InvalidParameter, AuthZ, AuthN, Internal
    message  string          // Developer message (English)
    metadata map[string]any  // Structured debugging info
    cause    error           // Wrapped original error
}

// Constructors for each error type
func NewNotFound(message string, cause error, metadata ...MetadataOption) *AppError { ... }
func NewInvalidParameter(message string, cause error, metadata ...MetadataOption) *AppError { ... }
func NewAuthZ(message string, cause error, metadata ...MetadataOption) *AppError { ... }

// Mapping to gRPC status codes
// NotFound          -> codes.NotFound
// InvalidParameter  -> codes.InvalidArgument
// AuthZ             -> codes.PermissionDenied
// AuthN             -> codes.Unauthenticated
// Internal          -> codes.Internal
```

Usage in handlers:
```go
if todoEnt == nil {
    return nil, errors.NewNotFound("todo not found", nil,
        errors.String("name", in.Name.String()),  // Metadata for debugging
    )
}
```

### When to Use
- Every domain error uses `errors.New*()` constructors
- Wrap errors with context: `fmt.Errorf("get todo: %w", err)`

### When NOT to Use
- DO NOT use stdlib's `errors.New("...")` for domain errors
- DO NOT use `status.Errorf()` directly — use AppError for consistent error format

---

## Exercises

### EX1: Trace Full Request Flow
Trace the complete processing flow of the `GetTodo` API from gRPC entry point to database query. Draw a sequence diagram showing how the request passes through each Clean Architecture layer.

**Files to read** (in order):
1. `todos/internal/handler/grpc/service/todos/get_todo.go` (Handler layer)
2. `todos/internal/usecase/todos/todo.go` (UseCase interface)
3. `todos/internal/service/todos/todo_getter.go` (Service implementation)
4. `todos/internal/domain/gateway/todo.go` (Gateway interface)
5. `todos/internal/infrastructure/datastore/todo_reader.go` (Infrastructure implementation)

**Deliverable**: A markdown document with sequence diagram and explanation of each layer's role in the request flow.

### EX2: Map Directory Structure and Identify Layers
Map the complete directory structure of the `todos/` project. For each directory:
- Identify which Clean Architecture layer it belongs to
- List all gateway interfaces in `domain/gateway/` and explain each interface's purpose
- For each file, explain why it belongs to that layer

**Deliverable**: A markdown document with annotated directory tree and gateway interface catalog.

### EX3: Local Dev Setup + List RPC Methods
1. Set up local development environment, successfully run `make run` and `make test-local`
2. Read the proto definition of `TodosService` and list all RPC methods with:
   - Method name
   - Request/Response messages
   - Description of functionality

**Deliverable**: Screenshot of successful `make run` + markdown document listing all RPC methods.

### EX4: Implement Mock Handler with 5-Step Pattern
Write a unit test for the `GetTodo` handler following the 5-step handler pattern (Parse → Build Input → Validate → Execute → Map Response). The test should:
- Mock the usecase interface
- Test happy path and error cases
- Demonstrate understanding of the interceptor chain (explain which interceptors run before your handler)

**Deliverable**: Working test file with at least 3 test cases + written explanation of interceptor execution order.
