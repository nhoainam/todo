# Phase 2: Database & Dependency Injection (Days 4-6)

## Goals
- Understand GORM ORM patterns in the project
- Master how Gateway interfaces are implemented
- Master Google Wire dependency injection

---

## 1. GORM ORM

### Problem Solved
Writing raw SQL:
- Error-prone: typos, SQL injection, wrong column names
- No compile-time checks
- Mapping between SQL rows and Go structs is manual and repetitive
- Transaction management is complex

**Real-world example:**
```go
// Raw SQL — error-prone
func GetTodo(db *sql.DB, id int64) (*Todo, error) {
    row := db.QueryRow("SELECT id, titl, status, priority FROM todos WHERE id = $1", id)
    //                          ^^^^ typo "titl" instead of "title" — only discovered at runtime!

    var todo Todo
    err := row.Scan(&todo.ID, &todo.Title, &todo.Status, &todo.Priority)
    // If you add a new column to SELECT → must add &todo.NewField to Scan (in correct order!)
    // Forget one field → runtime error: "Scan error: expected 5 destination arguments, got 4"

    return &todo, err
}

// SQL injection risk:
query := fmt.Sprintf("SELECT * FROM todos WHERE title = '%s'", userInput)
// userInput = "'; DROP TABLE todos; --" → deletes all data!
```

```go
// GORM — type-safe, auto-mapping
func (r *todoReader) Get(ctx context.Context, id entity.TodoID) (*entity.Todo, error) {
    todo, err := query.Use(db).Todo.WithContext(ctx).
        Where(query.Todo.ID.Eq(int64(id))).  // Type-safe, no typos
        First()                                // Auto-maps to struct
    return todo, err
}
// Adding a new column? Add field to struct → GORM maps it automatically
// SQL injection? GORM uses parameterized queries → safe
```

### Core Concept
GORM is an ORM (Object-Relational Mapping) for Go. It maps Go structs <-> database tables.

#### Entity as GORM Model

```go
// File: todos/internal/domain/entity/todo.go

type Todo struct {
    ID          TodoID
    TodoListID  *TodoListID
    Title       string
    Description *string
    Status      TodoStatus
    Priority    Priority
    DueDate     *time.Time
    AssigneeID  *UserID
    CreatorID   UserID

    // Relationships
    TodoList *TodoList `gorm:"foreignKey:TodoListID"`
    Tags     []*Tag    `gorm:"many2many:todo_tags;"`

    // Timestamps
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt  // Soft delete: DELETE -> UPDATE deleted_at
}

// Table name mapping
func (t *Todo) TableName() string { return "todos" }
```

**Key GORM features in the project:**
- **Soft delete**: `gorm.DeletedAt` — `DELETE` only sets `deleted_at`, doesn't remove the record
- **Preloading**: Load associations (e.g., `TodoList`, `Tags`)
- **AutoMigrate**: DO NOT use — the project uses Ridgepole

#### GORM Gen — Type-safe Queries

GORM Gen generates type-safe query builders:

```go
// Generated code (DO NOT edit directly)
query := query.Use(db)
todoQuery := query.Todo

// Type-safe queries
todo, err := todoQuery.WithContext(ctx).
    Where(todoQuery.ID.Eq(int64(todoID))).
    First()

// Type-safe joins
todoQuery.WithContext(ctx).
    Select(todoQuery.ALL).
    Join(todoTagQuery, todoQuery.ID.EqCol(todoTagQuery.TodoID)).
    Where(todoTagQuery.TagID.Eq(tagID.Int64())).
    Find()
```

### When to Use
- All database operations go through GORM
- Use GORM Gen for type-safe queries (preferred over raw GORM)
- Use Preload for associations when needed

### When NOT to Use
- DO NOT use `db.AutoMigrate()` — schema is managed by Ridgepole
- DO NOT use raw SQL unless GORM doesn't support it (extremely rare)
- DO NOT use GORM directly in handler/usecase — only in the infrastructure layer

---

## 2. Gateway Implementation Pattern

### Problem Solved
The domain layer needs database access but must not depend on GORM. An intermediate layer is needed.

**Real-world example:**
```go
// WITHOUT gateway — usecase calls GORM directly
func (s *todoGetter) Get(ctx context.Context, id int64) (*Todo, error) {
    var todo Todo
    err := s.db.Where("id = ?", id).First(&todo).Error  // Usecase knows about GORM!
    return &todo, err
}

// Problem 1: Unit tests must set up GORM + database
// Problem 2: Want to switch from GORM to sqlx? Must modify all usecases
// Problem 3: Usecase "knows too much" about infrastructure
```

```go
// WITH gateway — usecase only knows the interface
func (s *todoGetter) Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error) {
    todo, err := s.todoQueriesGateway.Get(ctx, in.ID, nil)  // Only calls the interface
    return &output.TodoGetter{Todo: todo}, err
}
// Unit test: mock the gateway, no database needed
// Switch databases: just implement a new gateway, usecases don't change
```

### Core Concept
Gateway interfaces (domain) + implementations (infrastructure):

```go
// INTERFACE — Domain layer
// File: todos/internal/domain/gateway/todo.go

// Separate Commands (write) and Queries (read)
type TodoCommandsGateway interface {
    Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
    Delete(ctx context.Context, todoID entity.TodoID) error
    Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
}

type TodoQueriesGateway interface {
    Get(ctx context.Context, todoID entity.TodoID, opts *GetTodoOptions) (*entity.Todo, error)
    List(ctx context.Context, opts *ListTodosOptions) (*query.OffsetPageResult[*entity.Todo], error)
    BatchGet(ctx context.Context, todoIDs []entity.TodoID) ([]*entity.Todo, error)
}
```

```go
// IMPLEMENTATION — Infrastructure layer
// File: todos/internal/infrastructure/datastore/todo_reader.go

type todoReader struct{}

func NewTodoReader() gateway.TodoQueriesGateway {
    return &todoReader{}
}

func (r *todoReader) Get(
    ctx context.Context,
    todoID entity.TodoID,
    opts *gateway.GetTodoOptions,
) (*entity.Todo, error) {
    // 1. Get DB connection from context
    db, err := DBFromContext(ctx)
    if err != nil { return nil, fmt.Errorf("get db from context: %w", err) }

    query := query.Use(db).Todo
    queryBuilder := query.WithContext(ctx)

    // 2. Apply options (preloading)
    if opts != nil && opts.With != nil {
        if _, ok := opts.With[entity.TodoAssociationTags]; ok {
            queryBuilder = queryBuilder.Preload(query.Tags)
        }
    }

    // 3. Execute query
    todo, err := queryBuilder.Where(query.ID.Eq(int64(todoID))).First()
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil  // Not found -> return nil (NOT an error)
        }
        return nil, fmt.Errorf("get todo: %w", err)
    }

    return todo, nil
}
```

#### Filter Pattern
```go
// Filter struct — Domain layer
type TodoFilter struct {
    StatusEq       *entity.TodoStatus
    PriorityEq     *entity.Priority
    DueDateGTE     *time.Time
    DueDateLTE     *time.Time
    TitleContain   *string
    AssigneeIDEq   *entity.UserID
    TodoListIDEq   *entity.TodoListID
}

// Apply filter — Infrastructure layer
func (r *todoReader) applyFilter(q *query.Query, builder query.ITodoDo, filter *TodoFilter) query.ITodoDo {
    if filter == nil { return builder }

    todoQuery := q.Todo
    if filter.StatusEq != nil {
        builder = builder.Where(todoQuery.Status.Eq(filter.StatusEq.Int()))
    }
    if filter.TitleContain != nil {
        builder = builder.Where(todoQuery.Title.Like("%" + *filter.TitleContain + "%"))
    }
    if filter.DueDateGTE != nil {
        builder = builder.Where(todoQuery.DueDate.Gte(*filter.DueDateGTE))
    }
    // ...
    return builder
}
```

#### Pagination Pattern
```go
// Domain query helpers
type OffsetPage struct {
    Offset int
    Limit  int
}

type OffsetPageResult[T any] struct {
    Items      []T
    TotalCount int64
    Page       *OffsetPage
}
```

### When to Use
- All database operations go through gateways
- Separate Commands/Queries when an entity has many operations
- Use Filter structs when queries have multiple optional conditions

### When NOT to Use
- DO NOT return gorm.DB outside of the infrastructure layer
- DO NOT put query logic in the service/usecase layer
- DO NOT merge Commands and Queries gateways for large entities

---

## 3. Database Context Pattern

### Problem Solved
Passing `*gorm.DB` directly through function parameters creates coupling. It's also hard to use transactions across multiple gateways.

**Real-world example:**
```go
// Passing DB via parameters — hard to use transactions
func (s *service) CreateTodoWithTags(db *gorm.DB, todo *Todo, tags []*Tag) error {
    // Create todo
    todoGW := NewTodoWriter(db)       // Must pass db to every gateway
    todoGW.Create(todo)

    // Create tags — IF this fails, todo was already created but tags weren't!
    tagGW := NewTodoTagWriter(db)     // Must pass db to every gateway
    tagGW.BatchCreate(tags)           // If error → data inconsistency!

    // Want to use a transaction? Must pass tx instead of db:
    tx := db.Begin()
    todoGW := NewTodoWriter(tx)       // Must recreate all gateways with tx
    tagGW := NewTodoTagWriter(tx)     // Recreate...
    // Very complex and easy to forget!
}
```

```go
// DB in context — gateways automatically share the same connection/transaction
func (s *service) CreateTodoWithTags(ctx context.Context, ...) error {
    ctx = s.binder.Bind(ctx)  // Create transaction, store in context
    // All gateways automatically use the SAME transaction from context
    s.todoGateway.Create(ctx, todo)     // Uses tx from context
    s.tagGateway.BatchCreate(ctx, tags) // Uses the SAME tx from context
    // If either fails → both rollback
}
```

### Core Concept
The DB connection is stored in context and extracted when needed:

```go
// Infrastructure: inject DB into context
ctx = gorm.WithContext(ctx, db)

// Gateway: extract DB from context
db, err := DBFromContext(ctx)
```

Transaction:
```go
// Binder gateway creates a transaction and injects it into context
func (b *binder) Bind(ctx context.Context) context.Context {
    tx := b.db.Begin()
    return gorm.WithContext(ctx, tx)
}

// Usecase uses the binder
func (g *todoGetter) Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error) {
    ctx = g.binder.Bind(ctx)  // Start transaction
    // ... all gateway calls in this context use the same transaction
}
```

### When to Use
- When multiple gateways need to run in the same transaction
- When using read/write splitting (context can contain a read-only DB)

### When NOT to Use
- Simple read-only operations may not need a transaction

---

## 4. Google Wire — Dependency Injection

### Problem Solved
Manually wiring dependencies:

**Real-world example:**
```go
// Manual wiring — 50+ dependencies
func main() {
    db := gorm.Open(...)
    todoReader := datastore.NewTodoReader()
    todoWriter := datastore.NewTodoWriter()
    tagWriter := datastore.NewTagWriter()
    todoHelper := helper.NewTodoHelper(todoReader)
    todoGetter := service.NewTodoGetter(todoReader, todoHelper, clock, binder, cfg)
    todoDeleter := service.NewTodoDeleter(todoWriter, tagWriter, todoHelper, binder)
    // ... 30 more lines for other dependencies
    todosService := handler.NewService(todoGetter, todoDeleter, ...) // 20+ params
    server := grpc.NewServer(todosService, adminService, ...)
}

// Problem 1: Adding a new dependency to NewTodoGetter?
//   → Must modify main.go, create the dependency, pass it in
// Problem 2: Forgot to pass a dependency? Runtime panic, not a compile error
// Problem 3: Circular dependency? Only discovered at runtime
// Problem 4: 200 lines of wiring code in main.go — who maintains this?
```

### Core Concept
Wire is compile-time dependency injection. Code generation, NOT runtime reflection.

**3 key concepts:**

#### Provider — Function that creates a dependency
```go
// Every constructor is a provider
func NewTodoReader() gateway.TodoQueriesGateway {
    return &todoReader{}
}

func NewTodoGetter(
    todoQueriesGateway gateway.TodoQueriesGateway,
    // ... Wire automatically resolves these dependencies
) usecase.TodoGetter {
    return &todoGetter{...}
}
```

#### WireSet — Group providers by package
```go
// File: todos/internal/infrastructure/datastore/wire.go
var WireSet = wire.NewSet(
    NewTodoReader,     // -> gateway.TodoQueriesGateway
    NewTodoWriter,     // -> gateway.TodoCommandsGateway
    NewTagWriter,      // -> gateway.TagCommandsGateway
    // ...
)

// File: todos/internal/service/service.go
var WireSet = wire.NewSet(
    todos.WireSet,
    admin.WireSet,
    helper.WireSet,
)

// File: todos/internal/infrastructure/infrastructure.go
var WireSet = wire.NewSet(
    datastore.WireSet,
    gorm.WireSet,
)
```

#### Injector — Composition root
```go
// File: todos/internal/registry/wire.go

//go:build wireinject

func InitializeServer(cfg *config.Config, logger *zap.Logger) (*grpc.Server, func(), error) {
    wire.Build(
        handler.WireSet,         // gRPC handlers, interceptors
        infrastructure.WireSet,  // GORM, external clients
        service.WireSet,         // UseCase implementations
        utils.WireSet,           // Utilities (clock, ID gen, validator)
    )
    return nil, nil, nil  // Wire generates actual code
}
```

Run `go tool wire ./internal/registry` -> generates `wire_gen.go`:
```go
// wire_gen.go (auto-generated — DO NOT EDIT)
func InitializeServer(cfg *config.Config, logger *zap.Logger) (*grpc.Server, func(), error) {
    db, cleanup, err := gorm.Open(cfg, logger)
    todoReader := datastore.NewTodoReader()
    todoWriter := datastore.NewTodoWriter()
    // ... Wire automatically wires 100+ dependencies
    server, cleanup2, err := grpc.NewServer(todosService, adminService, ...)
    return server, func() { cleanup2(); cleanup() }, nil
}
```

### When to Use
- All dependency wiring goes through Wire
- Every new package needs a WireSet
- Constructor functions should return interface types (so Wire can match them)

### When NOT to Use
- DO NOT use Wire for simple utility functions
- DO NOT create WireSets for functions that aren't constructors
- DO NOT edit `wire_gen.go` directly — run `go tool wire` to regenerate

---

## 5. Database Schema Management (Ridgepole)

### Problem Solved
Traditional schema migration tools (up/down migrations):
- Conflicts when multiple branches change the schema simultaneously
- Hard to know the current state of the schema
- Rollbacks are risky

**Real-world example:**
```
// Traditional migrations — up/down files
migrations/
  001_create_todos.up.sql          -- CREATE TABLE todos ...
  001_create_todos.down.sql        -- DROP TABLE todos
  002_add_priority_column.up.sql   -- ALTER TABLE todos ADD COLUMN priority ...
  002_add_priority_column.down.sql -- ALTER TABLE todos DROP COLUMN priority
  003_add_due_date.up.sql          -- ALTER TABLE todos ADD COLUMN due_date ...

// Problem 1: Branch A adds column "priority" (migration 002)
//            Branch B adds column "assignee_id" (also migration 002)
//            Merge → conflict! Must renumber migrations

// Problem 2: Production is at migration 002. What does the schema look like?
//            Must read all migrations from 001 → 002 to understand

// Problem 3: Rollback migration 002? down.sql may not be safe
//            (e.g.: DROP COLUMN loses data)
```

```ruby
// Ridgepole — declarative, just define the desired state
# todos.schema — THIS is the current state of the schema
create_table "todos" do |t|
  t.string  "title", null: false
  t.integer "priority", default: 0     # Branch A adds
  t.bigint  "assignee_id"              # Branch B adds
  t.datetime "due_date"
end
# Ridgepole automatically calculates needed ALTER statements
# Merge conflict? Just a text conflict in one file — easy to resolve
```

### Core Concept
Ridgepole is declarative schema management — you define the desired state, Ridgepole automatically generates ALTER statements.

```ruby
# File: todos/database/schemas/todos.schema

create_table "todos", id: :bigint, force: :cascade do |t|
  t.bigint  "todo_list_id"
  t.string  "title", null: false
  t.text    "description"
  t.integer "status", null: false, default: 0
  t.integer "priority", null: false, default: 0
  t.datetime "due_date"
  t.bigint  "assignee_id"
  t.bigint  "creator_id", null: false
  t.timestamps
  t.datetime "deleted_at"
end
```

```bash
# Apply schema changes
make migration
```

### When to Use
- All database schema changes update the `.schema` file
- Run `make migration` to apply

### When NOT to Use
- DO NOT use GORM AutoMigrate
- DO NOT run raw ALTER TABLE directly

---

## Exercises

### EX1: Gateway Tracing
Read the `TodoQueriesGateway` interface and find its corresponding implementation.
Map each interface method to its corresponding GORM query.

### EX2: Write GORM Query
Write a new GORM query: find todos by `priority` and `due_date` range.
- Add filter fields to `TodoFilter`
- Implement filter logic in `applyFilter`

### EX3: Wire Dependency Graph
Trace the Wire dependency graph: starting from `InitializeServer`, draw the tree of dependencies.
Identify what providers each WireSet provides.

### EX4: Add New Provider
Add a new helper service to the Wire graph:
1. Create an interface in the domain layer
2. Create the implementation
3. Add it to the WireSet
4. Run `go tool wire` and verify `wire_gen.go`
