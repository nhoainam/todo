# Phase 5: Observability & End-to-End Feature (Days 11-13)

## Goals
- Understand the observability stack (Datadog, Sentry, Zap)
- Master the code generation workflow
- Practice implementing a feature end-to-end

---

## 1. Structured Logging with Zap

### Problem Solved
`fmt.Println` or `log.Printf`:
- No log levels (debug, info, warn, error)
- No structured data — hard to parse/search in log aggregation tools
- No context (request ID, user ID, trace ID)
- Slow (reflection-based formatting)

**Real-world example:**
```go
// fmt.Println / log.Printf — production nightmare
fmt.Println("todo created successfully")
log.Printf("error getting todo: %v", err)
log.Printf("processing todo %d for user %d, priority: %d", todoID, userID, priority)

// In Datadog/Kibana, logs look like this:
// 2025-01-01 10:00:01 todo created successfully
// 2025-01-01 10:00:02 error getting todo: connection refused
// 2025-01-01 10:00:03 processing todo 456 for user 123, priority: 2

// Problem 1: Search "all errors" → can't because there's no log level
// Problem 2: Search "all logs for user 123" → must regex match text
// Problem 3: Production has 10000 logs/sec → want to disable debug logs? Can't
// Problem 4: "connection refused" error relates to which request? Unknown
```

```go
// Zap structured logging — easy to search, filter, aggregate
logger.Info("todo created",
    zap.Int64("todo_id", 456),
    zap.Int64("user_id", 123),
)
// Output JSON: {"level":"info","msg":"todo created","todo_id":456,"user_id":123}
// Search user_id=123 → exact, no false positives
// Filter level=error → only shows errors
// Production → set level=warn → automatically disables info/debug logs
```

### Core Concept
Zap is a structured, leveled, high-performance logger.

```go
// Structured logging — each field is a key-value pair
logger.Info("todo created",
    zap.String("todo_id", todo.ID.String()),
    zap.String("user_id", userID.String()),
    zap.String("status", todo.Status.String()),
    zap.Duration("processing_time", elapsed),
)
// Output (JSON): {"level":"info","msg":"todo created","todo_id":"456","user_id":"123","status":"pending","processing_time":"150ms"}

// Error logging with error object
logger.Error("failed to get todo",
    zap.Error(err),
    zap.String("todo_id", todoID.String()),
)

// Context-aware logging (includes request metadata)
log.InfoContext(ctx, "processing request",
    zap.String("method", "GetTodo"),
)
```

#### Log Levels
```
Debug  → Development only, verbose details
Info   → Normal operations (request received, todo created)
Warn   → Something unexpected but handled (retry, fallback)
Error  → Something failed, needs attention
Fatal  → Application cannot continue, exits
```

#### Configuration

```go
func New(cfg *config.Config) (*zap.Logger, error) {
    if cfg == nil || cfg.Environment.IsLocal() {
        return zap.NewDevelopment()  // Human-readable
    }
    return zap.NewProduction()       // JSON format for log aggregation
}
```

### When to Use
- Every significant operation (API calls, DB queries, external service calls)
- Error cases MUST always be logged
- Use structured fields for searchability

### When NOT to Use
- DO NOT log sensitive data (passwords, tokens, PII)
- DO NOT log in hot paths unnecessarily (inner loops)
- DO NOT use `fmt.Println` — always use the logger

---

## 2. Distributed Tracing with Datadog

### Problem Solved
Microservice architecture — 1 request goes through multiple services. When there's an error:
- Which service is slow?
- What path did the request take?
- Where is the bottleneck?

**Real-world example:**
```
// User report: "The todo detail page takes 5 seconds to load!"

// WITHOUT tracing — how do you debug?
// Check BFF logs: "GetTodo took 4.8s" — but where's the slowness?
// Check Todos service logs: "GetTodo took 0.1s" — that's fast?
// Check DB logs: "query took 50ms" — DB is fast too?
// 4.8s - 0.1s = 4.7s where did it go??? Network? Serialization? What?
// Must manually correlate logs across services by timestamp — takes hours

// WITH tracing — 1 click to see everything:
// Datadog Trace View:
// [BFF: GetTodo]──────────────────────── 4.8s
//   ├─[gRPC call: Todos.GetTodo]──────── 0.1s
//   ├─[gRPC call: Users.GetUser]──────── 4.5s  ← BOTTLENECK!
//   └─[Response mapping]──────────────── 0.2s
// Clear: Users.GetUser is slow at 4.5s → fix in Users service
```

### Core Concept
Datadog APM traces every request across services. Auto-instrumented via interceptors/middleware.

```go
func Init(cfg *config.Config, logger *zap.Logger) (func(), error) {
    tracer.Start(
        tracer.WithService(cfg.ServiceName),
        tracer.WithEnv(string(cfg.Environment)),
    )
    return tracer.Stop, nil
}
```

#### Auto-instrumentation

```go
// gRPC server — every RPC call is traced
grpctrace.UnaryServerInterceptor(
    grpctrace.WithServiceName("andpad-todos"),
)

// gRPC client (BFF -> Todos) — every call is traced
grpctrace.UnaryClientInterceptor(
    grpctrace.WithServiceName("andpad-todos-bff"),
)

// HTTP server (BFF)
chitrace.Middleware(
    chitrace.WithServiceName("andpad-todos-bff"),
)

// GORM — every DB query is traced
gormtrace.Open(...)
```

**Trace flow:**
```
[Browser] → [BFF: HTTP trace] → [BFF: gRPC client trace] → [Todos: gRPC server trace] → [Todos: GORM DB trace]
```

Each trace has:
- **Trace ID**: Unique per request, propagated across services
- **Span**: Unit of work (1 gRPC call, 1 DB query)
- **Tags**: Metadata (service name, resource, error)
- **Duration**: Latency of each span

### When to Use
- All service boundaries (gRPC, HTTP)
- Database queries
- External API calls
- Significant internal operations

### When NOT to Use
- Tracing is already auto-instrumented — no need for manual tracing
- Trivial operations (string formatting, in-memory operations)

---

## 3. Error Reporting with Sentry

### Problem Solved
Errors get lost in logs. Need:
- Aggregation: group similar errors into 1 issue
- Alerting: notify when a new error occurs
- Context: user info, request data, stack trace
- Tracking: has the error been fixed?

**Real-world example:**
```
// ONLY using logs — errors get lost
// Production has 50000 log entries/day
// 200 of them are error logs — but who reads them?

// Log:
// 10:00:01 ERROR failed to get todo: connection refused
// 10:00:02 INFO  todo created successfully
// 10:00:03 ERROR failed to get todo: connection refused  (same as above!)
// ... 49997 more log entries ...
// 10:59:59 ERROR failed to update status: deadline exceeded

// Problem 1: 150/200 errors are "connection refused" — same bug but logged 150 times
// Problem 2: "deadline exceeded" is a new bug → nobody knows because it's buried in 50000 logs
// Problem 3: Has the "connection refused" bug been fixed? Must manually search logs
// Problem 4: Which users are affected? Unknown because logs lack user context
```

```
// WITH Sentry — errors are grouped, tracked, alerted
// Sentry Dashboard:
// Issue #1: "failed to get todo: connection refused" — 150 events, 45 users affected
//   → Status: Resolved (fixed)
//   → First seen: 10:00:01, Last seen: 10:45:00
// Issue #2: "failed to update status: deadline exceeded" — 3 events, 2 users affected  [NEW!]
//   → Alert → Slack notification: "New issue in production!"
//   → Stack trace + user ID + request data → easy to debug
```

### Core Concept

```go
// Initialize
sentry.Init(cfg.Sentry.DSN, cfg.Sentry.Environment, logger)
defer sentry.Flush(10 * time.Second)

// Report error (automatically includes stack trace, context)
sentry.CaptureException(ctx, err)

// Set user context (from authn interceptor)
sentrygo.SetUser(sentrygo.User{
    ID: userID.String(),
})
```

#### Auto-capture in interceptors

```go
// gRPC interceptor
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
        resp, err := handler(ctx, req)
        if err != nil {
            // Automatically report unhandled errors to Sentry
            sentry.CaptureException(ctx, err)
        }
        return resp, err
    }
}
```

### When to Use
- Unexpected errors (database down, external service failures)
- Panics (via recovery interceptor)
- Business-critical failures

### When NOT to Use
- Expected errors (not found, validation errors) — don't spam Sentry
- When the error is already handled properly and returns an appropriate gRPC status

---

## 4. Code Generation Workflow

### Problem Solved
Manual code for: proto stubs, mocks, Wire wiring, enum helpers, GraphQL resolvers — repetitive, error-prone, out-of-sync.

**Real-world example:**
```go
// Manual mock — must write by hand every time the interface changes
type MockTodoQueriesGateway struct{}

func (m *MockTodoQueriesGateway) Get(ctx context.Context, id entity.TodoID, opts *gateway.GetTodoOptions) (*entity.Todo, error) {
    return &entity.Todo{ID: id}, nil
}

func (m *MockTodoQueriesGateway) List(ctx context.Context, opts *gateway.ListTodosOptions) (*query.OffsetPageResult[*entity.Todo], error) {
    return &query.OffsetPageResult[*entity.Todo]{}, nil
}

// Problem 1: Interface adds method BatchGet() → must update mock
// Problem 2: Method signature changes → must update mock
// Problem 3: 20 interfaces × average 4 methods = 80 mock methods written by hand!
// Problem 4: Mock DOESN'T have EXPECT() → can't verify if a method was called

// Same for Wire: add new dependency → must update wire code manually
// Same for GraphQL: add new field → must write resolver skeleton manually
```

`make generate` creates everything automatically — just run 1 command.

### Core Concept

```bash
# Backend service
cd todos
make generate
# Equivalent to:
#   go generate ./...              # mockgen, enumer, wire
#   sed -i ... wire_gen.go         # Fix wire build tags

# BFF
cd todos-bff
make generate
# Equivalent to:
#   go tool gqlgen generate .      # GraphQL code gen
#   go generate ./...              # mockgen, wire
```

#### Types of code gen in the project:

| Tool | Input | Output | Trigger |
|---|---|---|---|
| **gqlgen** | `schema.graphqls` + `gqlgen.yml` | Resolvers, models, generated code | Schema changes |
| **Wire** | `wire.go` (injector definitions) | `wire_gen.go` (wiring code) | Dependency changes |
| **mockgen** | Interface with `//go:generate` | Mock implementation | Interface changes |
| **enumer** | Enum type with `//go:generate` | String(), validation methods | Enum values change |
| **protoc** | `.proto` files | Go stubs | Proto schema changes |

#### When to run code gen?

```
Change schema.graphqls     → make generate (in BFF)
Change gateway interface   → go generate ./internal/domain/gateway/...
Add new dependency         → go tool wire ./internal/registry
Add new enum value         → go generate ./internal/domain/entity/...
```

### When to Use
- AFTER changing source files (schema, interfaces, wire.go)
- BEFORE committing code

### When NOT to Use
- DO NOT edit generated files directly — they will be overwritten
- DO NOT commit generated files that haven't been regenerated

---

## 5. Comprehensive Exercise: Implement Feature End-to-End

**Feature**: Add "Add Tag to Todo" functionality (assign a tag to a todo item)

### Backend (Todos Service)

#### Step 1: Database Schema
```ruby
# File: todos/database/schemas/todo_tags.schema

create_table "todo_tags", id: :bigint, force: :cascade do |t|
  t.bigint  "todo_id", null: false
  t.bigint  "tag_id", null: false
  t.timestamps
  t.index ["todo_id", "tag_id"], unique: true
end
```

#### Step 2: Domain Entity
```go
// File: todos/internal/domain/entity/todo_tag.go

type TodoTag struct {
    ID        int64
    TodoID    TodoID
    TagID     TagID
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (tt *TodoTag) TableName() string { return "todo_tags" }
```

#### Step 3: Gateway Interfaces
```go
// File: todos/internal/domain/gateway/todo_tag.go

//go:generate go tool mockgen -destination=mock/$GOFILE -source=$GOFILE

type TodoTagCommandsGateway interface {
    Create(ctx context.Context, todoTag *entity.TodoTag) (*entity.TodoTag, error)
    Delete(ctx context.Context, todoID entity.TodoID, tagID entity.TagID) error
}

type TodoTagQueriesGateway interface {
    List(ctx context.Context, todoID entity.TodoID) ([]*entity.TodoTag, error)
    Exists(ctx context.Context, todoID entity.TodoID, tagID entity.TagID) (bool, error)
}
```

#### Step 4: Gateway Implementation
```go
// File: todos/internal/infrastructure/datastore/todo_tag_reader.go

type todoTagReader struct{}

func NewTodoTagReader() gateway.TodoTagQueriesGateway {
    return &todoTagReader{}
}

func (r *todoTagReader) Exists(ctx context.Context, todoID entity.TodoID, tagID entity.TagID) (bool, error) {
    db, err := DBFromContext(ctx)
    if err != nil { return false, err }
    // GORM query...
}
```

#### Step 5: UseCase Interfaces
```go
// File: todos/internal/usecase/todos/todo_tag.go

type TodoTagger interface {
    AddTag(ctx context.Context, in *input.TodoTagger) (*output.TodoTagger, error)
}

type TodoUntagger interface {
    RemoveTag(ctx context.Context, in *input.TodoUntagger) (*output.TodoUntagger, error)
}
```

#### Step 6: Service Implementation
```go
// File: todos/internal/service/todos/todo_tagger.go

type todoTagger struct {
    todoTagCommandsGateway gateway.TodoTagCommandsGateway
    todoHelper             sharedport.TodoHelper
    binder                 gateway.Binder
}

func (s *todoTagger) AddTag(ctx context.Context, in *input.TodoTagger) (*output.TodoTagger, error) {
    ctx = s.binder.Bind(ctx)

    // 1. Verify todo exists and user has access
    todo, err := s.todoHelper.Get(ctx, in.TodoName)
    if err != nil { return nil, err }

    // 2. Create tag association
    todoTag := &entity.TodoTag{
        TodoID: todo.ID,
        TagID:  in.TagID.ToEntity(),
    }
    todoTag, err = s.todoTagCommandsGateway.Create(ctx, todoTag)
    if err != nil { return nil, err }

    return &output.TodoTagger{TodoTag: todoTag}, nil
}
```

#### Step 7: gRPC Handler
```go
// File: todos/internal/handler/grpc/service/todos/add_tag.go

func (s *service) AddTag(ctx context.Context, req *todospb.AddTagRequest) (*todospb.AddTagResponse, error) {
    name, err := todos.ParseTodoResourceName(req.TodoName)
    if err != nil { return nil, errors.NewInvalidParameter("invalid name", err) }

    in := input.TodoTagger{TodoName: cast.Value(name), TagID: req.TagId}
    if err := s.validator.Struct(&in); err != nil { return nil, errors.NewInvalidParameter("invalid request", err) }

    _, err = s.todoTagger.AddTag(ctx, &in)
    if err != nil { return nil, err }

    return &todospb.AddTagResponse{}, nil
}
```

#### Step 8: Wire Registration
```go
// Add to corresponding WireSets
// datastore/wire.go: NewTodoTagReader, NewTodoTagWriter
// service/todos/wire.go: NewTodoTagger, NewTodoUntagger
```

#### Step 9: Tests
- Unit test for `todoTagger.AddTag()` with mocks
- Integration test for the `AddTag` gRPC endpoint

### Frontend Gateway (BFF)

#### Step 10: GraphQL Schema
```graphql
# Add to schema.graphqls
extend type Mutation {
    addTagToTodo(todoName: ResourceName!, tagId: Int64!): AddTagPayload!
    removeTagFromTodo(todoName: ResourceName!, tagId: Int64!): RemoveTagPayload!
}

type AddTagPayload {
    success: Boolean!
}
```

#### Step 11: Generate & Implement
```bash
make generate  # Generate resolver stubs
```

```go
// Implement resolver
func (r *mutationResolver) AddTagToTodo(ctx context.Context, todoName scalar.ResourceName, tagID int64) (*model.AddTagPayload, error) {
    // Call usecase -> gRPC client -> Todos service
}
```

---

## Phase 5 Checklist

- [ ] Understand Zap logging levels and structured fields
- [ ] Understand Datadog trace flow across services
- [ ] Understand Sentry error reporting flow
- [ ] Successfully run `make generate`
- [ ] Implement a feature end-to-end (backend + BFF)
- [ ] Write unit tests and integration tests for the new feature
