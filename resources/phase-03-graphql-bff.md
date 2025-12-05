# Phase 3: GraphQL & BFF Pattern (Days 7-8)

## Goals
- Understand the BFF (Backend For Frontend) pattern
- Master gqlgen code generation and resolver pattern
- Understand the DataLoader pattern for N+1 prevention

---

## 1. BFF (Backend For Frontend) Pattern

### Problem Solved
Frontend calling multiple backend microservices directly:
- Frontend must know the address of each service
- Multiple round-trips to get enough data for one page
- Backend APIs designed for general use, not optimized for UI
- Auth logic duplicated on the frontend
- CORS issues with multiple origins

**Real-world example:**
```
// WITHOUT BFF — Frontend calls multiple services directly
// "Todo Detail" page needs:

// Request 1: Get todo info
fetch("https://todos-api.internal/todos/456")       // Todos service

// Request 2: Get creator info
fetch("https://users-api.internal/users/789")       // Users service

// Request 3: Get comments
fetch("https://todos-api.internal/todos/456/comments")  // Todos service

// Problem 1: 3 round-trips → slow, especially on mobile
// Problem 2: Frontend knows 3 service URLs → tight coupling
// Problem 3: Auth token must be sent to multiple services → security risk
// Problem 4: Todos service changes URL? Must update frontend
```

```
// WITH BFF — Frontend calls only 1 endpoint
// Request: 1 GraphQL query → BFF aggregates internally
fetch("https://bff.api/query", {
    body: `{ todo(name: "users/1/todo-lists/2/todos/456") {
        title, creator { name }, comments { ... }
    }}`
})
// BFF internally calls services via gRPC (fast, binary)
// Frontend only knows 1 URL, 1 auth flow
```

### Core Concept
BFF is a gateway layer between frontend and backend services:

```
[Frontend]          [Todos-BFF]                [Backend Services]
    │                    │                           │
    │  GraphQL query     │                           │
    │───────────────────>│                           │
    │                    │  gRPC: GetTodo             │
    │                    │─────────────────────────> │ Todos
    │                    │  gRPC: GetUser             │
    │                    │─────────────────────────> │ Users
    │                    │                           │
    │  Single JSON resp  │                           │
    │<───────────────────│                           │
```

**BFF aggregates data from backend services:**
- `Todos` — Todo/TodoList management
- `Users` — User/Organization management

### Architecture
```go
// File: todos-bff/cmd/todos-bff/main.go

func main() {
    cfg, _ := config.Load()
    logger, _ := zaputil.New(cfg)

    // Initialize all dependencies via Wire
    httpServer, cleanup, _ := registry.InitHttpServer(cfg, logger)
    defer cleanup()

    // HTTP server (not gRPC — BFF serves GraphQL over HTTP)
    server := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
        Handler:      httpServer,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 30 * time.Second,
    }

    // Graceful shutdown...
}
```

### When to Use
- Frontend needs data from multiple backend services in 1 request
- Frontend needs a flexible query language (GraphQL)
- Need centralized auth, CORS, rate limiting for the frontend

### When NOT to Use
- Service-to-service communication (use gRPC directly)
- Simple CRUD without aggregation

---

## 2. GraphQL with gqlgen

### Problem Solved
REST API:
- Over-fetching: API returns more data than needed
- Under-fetching: Need to call multiple endpoints to get enough data
- Schema evolution is hard: adding new fields can break old clients
- No introspection — clients don't know what the API offers

**Real-world example:**
```
// REST — over-fetching and under-fetching
// Todo list page only needs: title, status
GET /api/todos
Response: [
  {
    "id": 1, "title": "Buy groceries", "status": "pending",
    "description": "Milk, eggs, bread...", "priority": "medium",
    "assignee_id": 789, "created_at": "...", "updated_at": "...",
    "due_date": "2025-03-15", "todo_list_id": 5,
    // ... 10 more fields the frontend DOESN'T NEED → wasted bandwidth
  }
]

// Todo detail page needs: todo + creator name
GET /api/todos/1              // Request 1: get todo
GET /api/users/789            // Request 2: get creator — under-fetching!
// 2 requests for 1 page!
```

```graphql
# GraphQL — fetch exactly what you need
# List page:
query { todos { nodes { title, status } } }

# Detail page:
query { todo(name: "...") { title, status, creator { name } } }
# 1 request, only returns requested fields
```

### Core Concept
GraphQL lets clients specify exactly what data they need.

#### gqlgen Workflow

```
1. Define schema (.graphqls)
       ↓
2. Configure (gqlgen.yml)
       ↓
3. Generate code (make generate)
       ↓
4. Implement resolvers
```

**Step 1: Schema** — `graph/schema.graphqls`

```graphql
# Custom scalars
scalar Time
scalar ResourceName

# Directives
directive @hasPermission(permissions: [Permission!]!) on FIELD_DEFINITION
directive @validateInput on FIELD_DEFINITION

# Types
type Todo {
  name: ResourceName!
  title: String!
  description: String
  status: TodoStatus!
  priority: Priority!
  dueDate: Time
  createdAt: Time!

  # Field resolvers (loaded on demand)
  creator: User
  assignee: User
  todoList: TodoList
  tags: [Tag!]!
}

# Queries
type Query {
  todos(listName: ResourceName!, input: TodosInput): TodoConnection!
  todo(name: ResourceName!): Todo!
}

# Mutations
type Mutation {
  createTodo(input: CreateTodoInput!): CreateTodoPayload!
  deleteTodo(name: ResourceName!): DeleteTodoPayload!
  updateTodoStatus(name: ResourceName!, status: TodoStatus!): UpdateTodoStatusPayload!
}
```

**Step 2: Configuration** — `gqlgen.yml`

```yaml
schema:
  - graph/*.graphqls

exec:
  filename: internal/handler/graph/generated/generated.go
  package: generated

model:
  filename: internal/handler/graph/model/models_gen.go
  package: model

resolver:
  layout: follow-schema
  dir: internal/handler/graph
  package: graph

# Custom type mappings
models:
  Time:
    model: .../scalar.Time
  ResourceName:
    model: .../scalar.ResourceName

  # Field-level resolvers
  Todo:
    fields:
      creator:
        resolver: true    # Generate resolver function for this field
      assignee:
        resolver: true
      todoList:
        resolver: true
      tags:
        resolver: true
```

**Step 3: Generate** — `make generate`

```bash
go tool gqlgen generate .   # Generate resolvers, models, generated code
go generate ./...           # Generate wire, mocks
```

**Step 4: Implement resolvers**

### When to Use
- Frontend-facing APIs
- When clients need flexible data fetching
- When aggregating data from multiple sources

### When NOT to Use
- File upload/download (use REST endpoints)
- Service-to-service (use gRPC)
- Simple APIs with 1-2 endpoints (overkill)

---

## 3. Resolver Pattern

### Problem Solved
The GraphQL schema defines the data structure, but who fetches the data and from where?

**Real-world example:**
```graphql
# Schema defines "what the data looks like"
type Todo {
  title: String!
  creator: User       # Where does User come from? Which database? Which service?
  todoList: TodoList   # How is TodoList fetched? One query or a join?
  tags: [Tag!]!
}
# Schema DOESN'T say: "fetch creator from Users service via gRPC"
# Schema DOESN'T say: "fetch todoList by querying todo_lists table with todo_list_id"
# → Resolvers are needed to "resolve" how to fetch data for each field
```

### Core Concept
Resolvers are functions that gqlgen calls to resolve each field in the schema.

#### Root Resolver

```go
// File: todos-bff/internal/handler/graph/resolver.go

type Resolver struct {
    validator *validator.Validate

    // All usecase dependencies
    todosLister       usecase.TodosLister
    todoGetter        usecase.TodoGetter
    todoCreator       usecase.TodoCreator
    todoDeleter       usecase.TodoDeleter
    todoStatusUpdater usecase.TodoStatusUpdater
    // ... more usecases
}

func New(
    todosLister usecase.TodosLister,
    todoGetter usecase.TodoGetter,
    // ... all dependencies injected by Wire
    validator *validator.Validate,
) generated.Config {
    cfg := generated.Config{
        Resolvers: &Resolver{...},
    }
    cfg.Directives.ValidateInput = directives.ValidateInput(validator)
    cfg.Directives.HasPermission = directives.HasPermission()
    return cfg
}
```

#### Query Resolver
```go
// File: todos-bff/internal/handler/graph/*.resolvers.go

func (r *queryResolver) Todos(
    ctx context.Context,
    listName scalar.ResourceName,
    input *model.TodosInput,
) (*model.TodoConnection, error) {
    // 1. Map GraphQL input -> usecase input
    usecaseInput := mapper.TodoListInputFromGQL(listName, input)

    // 2. Validate
    if err := r.validator.Struct(usecaseInput); err != nil { ... }

    // 3. Call usecase
    out, err := r.todosLister.List(ctx, usecaseInput)
    if err != nil { return nil, err }

    // 4. Map output -> GraphQL model
    return mapper.TodoConnectionFromUsecaseOutput(out), nil
}
```

#### Field Resolver (Lazy Loading)
```go
// Called IF AND ONLY IF the client queries the "creator" field
func (r *todoResolver) Creator(
    ctx context.Context,
    obj *model.Todo,
) (*model.User, error) {
    // Uses DataLoader to batch load (see section 4)
    return dataloaders.For(ctx).UserLoader.Load(ctx, obj.CreatorName)()
}

// Called IF AND ONLY IF the client queries the "todoList" field
func (r *todoResolver) TodoList(
    ctx context.Context,
    obj *model.Todo,
) (*model.TodoList, error) {
    return dataloaders.For(ctx).TodoListLoader.Load(ctx, obj.TodoListName)()
}
```

**When is a field resolver called?**

```graphql
# Query 1: Only fetches title and status — Creator resolver is NOT called
query {
  todo(name: "users/1/todo-lists/2/todos/3") {
    title
    status
  }
}

# Query 2: Includes creator — Creator resolver IS called
query {
  todo(name: "users/1/todo-lists/2/todos/3") {
    title
    creator {
      name
    }
  }
}
```

### When to Use
- Query/Mutation resolvers for root operations
- Field resolvers for relationships (creator, assignee, todoList, tags)
- Field resolvers for expensive computations (only computed when requested)

### When NOT to Use
- Simple scalar fields (title, status, createdAt) — gqlgen auto-resolves these

---

## 4. DataLoader Pattern

### Problem Solved
N+1 query problem:

**Real-world example:**
```graphql
# Client query: list 20 todos with their creators
query {
  todos(listName: "users/1/todo-lists/2") {
    nodes {
      title
      creator { name }    # <- Field resolver is called for EVERY todo
    }
  }
}
```

```
# WITHOUT DataLoader — each todo makes a separate request:
1. List 20 todos                               → 1 gRPC call
2. Get creator for todo 1 (user_id=10)         → 1 gRPC call
3. Get creator for todo 2 (user_id=10)         → 1 gRPC call (SAME user! but called again)
4. Get creator for todo 3 (user_id=20)         → 1 gRPC call
...
21. Get creator for todo 20 (user_id=15)       → 1 gRPC call
Total: 21 gRPC calls! (N+1 problem)

# If many todos share the same creator → still calls multiple times for the same user!
```

### Core Concept

DataLoader batches multiple requests within the same tick into a single batch request.

```
Without DataLoader:          With DataLoader:
getCreator(todo1)  → query    getCreator(todo1)  ─┐
getCreator(todo2)  → query    getCreator(todo2)  ─┤ wait 1ms
getCreator(todo3)  → query    getCreator(todo3)  ─┤
...                           ...                 ─┤
getCreator(todo20) → query    getCreator(todo20) ─┘→ batchGetUsers([10,10,20,...,15]) → 1 query
Total: 20 queries             Total: 1 query
```

```go
// File: todos-bff/internal/handler/graph/dataloaders/dataloader.go

type Loaders struct {
    TodoLoader     *dataloader.Loader[string, *model.Todo]
    UserLoader     *dataloader.Loader[string, *model.User]
    TodoListLoader *dataloader.Loader[string, *model.TodoList]
}

func NewLoaders(
    todoQueriesGateway gateway.TodoQueriesGateway,
    userQueriesGateway gateway.UserQueriesGateway,
) *Loaders {
    userLoader := &userBatchLoader{
        userQueriesGateway: userQueriesGateway,
    }

    return &Loaders{
        UserLoader: dataloader.NewBatchedLoader(
            userLoader.batchGetUsers,
            dataloader.WithCache(dataloader.NoCache[string, *model.User]{}),  // No caching
            dataloader.WithWait[string, *model.User](time.Millisecond),       // 1ms batching window
        ),
    }
}
```

Batch function:
```go
func (l *userBatchLoader) batchGetUsers(
    ctx context.Context,
    keys []string,
) []*dataloader.Result[*model.User] {
    // 1 gRPC call to fetch multiple users
    users, err := l.userQueriesGateway.BatchGet(ctx, keys)
    // Map results in the order of keys
}
```

**Per-request DataLoader** — Each HTTP request creates new Loaders:

```go
// File: middleware/http/dataloaders/
func Middleware(factory *dataloaders.Factory) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            loaders := factory.NewLoaders()  // Fresh loaders per request
            ctx := context.WithValue(r.Context(), loadersKey, loaders)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### When to Use
- Field resolvers for relationships (1-to-1, many-to-1)
- When multiple items need to load the same type of data
- Backend has batch APIs (BatchGet, BatchList)

### When NOT to Use
- Root queries (todos, todoLists) — only 1 call, no need to batch
- Fields that are always included (not lazy-loaded)

---

## 5. GraphQL Directives & Scalars

### Problem Solved (Directives)
Every resolver needs to check permissions and validate input. Can't copy-paste every time.

**Real-world example:**
```go
// WITHOUT directives — copy-paste permission checks into every resolver
func (r *mutationResolver) DeleteTodo(ctx context.Context, name string) (*model.Payload, error) {
    // Copy-paste permission check
    if !hasPermission(ctx, "TODO_DELETE") {
        return nil, errors.New("permission denied")
    }
    // Copy-paste validation
    if err := validate(name); err != nil {
        return nil, err
    }
    // ... business logic
}

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.Input) (*model.Payload, error) {
    // AGAIN copy-paste permission check
    if !hasPermission(ctx, "TODO_CREATE") { return nil, errors.New("permission denied") }
    // AGAIN copy-paste validation
    if err := validate(input); err != nil { return nil, err }
    // ...
}
// 20 mutations → 20 copy-pastes. Forget once → security hole!
```

```graphql
# WITH directives — declarative in the schema
type Mutation {
  deleteTodo(name: ResourceName!): Payload! @hasPermission(permissions: [TODO_DELETE]) @validateInput
  createTodo(input: Input!): Payload! @hasPermission(permissions: [TODO_CREATE]) @validateInput
}
# Permission and validation are automatic, impossible to forget!
```

### Core Concept (Directives)

```graphql
# Schema
directive @hasPermission(permissions: [Permission!]!) on FIELD_DEFINITION
directive @validateInput on FIELD_DEFINITION

type Mutation {
  deleteTodo(name: ResourceName!): DeleteTodoPayload!
    @hasPermission(permissions: [TODO_DELETE])
    @validateInput
}
```

```go
// File: todos-bff/internal/handler/graph/directives/

func HasPermission() func(ctx context.Context, obj any, next graphql.Resolver, permissions []model.Permission) (any, error) {
    return func(ctx context.Context, obj any, next graphql.Resolver, permissions []model.Permission) (any, error) {
        // Check user permissions
        if !hasRequiredPermissions(ctx, permissions) {
            return nil, errors.NewAuthZ("insufficient permissions")
        }
        return next(ctx)  // Continue to resolver
    }
}

func ValidateInput(v *validator.Validate) func(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
    return func(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
        // Validate input using go-playground/validator
        return next(ctx)
    }
}
```

### Problem Solved (Scalars)
GraphQL built-in scalars (String, Int, Boolean) aren't sufficient for domain-specific types.

**Real-world example:**
```graphql
# ONLY using built-in scalars
type Todo {
  name: String!       # "users/1/todo-lists/2/todos/3" — it's a String but has a specific format
  createdAt: String!  # "2025-01-01T00:00:00Z" — it's a String but must be ISO 8601
  todoListId: Int!    # 1073741824 — Int overflow! GraphQL Int is 32-bit
}
# Problem: Client sends name = "hello world" → not a resource name but still a valid String
# Problem: Client sends createdAt = "tomorrow" → not a datetime but still a valid String
# Validation must be done manually in every resolver
```

```graphql
# WITH custom scalars — automatic validation
scalar ResourceName  # Automatically validates format "users/{id}/todo-lists/{id}/todos/{id}"
scalar Time          # Automatically parses/validates ISO 8601 datetime
scalar Int64         # Supports 64-bit integers

type Todo {
  name: ResourceName!    # Client sends wrong format → error before reaching the resolver
  createdAt: Time!       # Automatically serializes/deserializes
}
```

### Core Concept (Scalars)

```go
// File: todos-bff/internal/handler/graph/scalar/

type ResourceName string

// MarshalGQL serializes ResourceName for GraphQL responses
func (r ResourceName) MarshalGQL(w io.Writer) {
    io.WriteString(w, strconv.Quote(string(r)))
}

// UnmarshalGQL deserializes ResourceName from GraphQL input
func (r *ResourceName) UnmarshalGQL(v any) error {
    str, ok := v.(string)
    if !ok { return fmt.Errorf("resource name must be a string") }
    *r = ResourceName(str)
    return nil
}
```

### When to Use (Directives)
- Cross-cutting concerns across many resolvers (auth, validation)
- Declarative policies in the schema

### When NOT to Use (Directives)
- Business logic specific to a single resolver

---

## 6. HTTP Middleware (BFF)

### Problem Solved
BFF serves HTTP (not gRPC). Needs middleware for auth, CORS, logging, similar to gRPC interceptors.

**Real-world example:**
```go
// WITHOUT middleware — every HTTP handler handles everything itself
func graphqlHandler(w http.ResponseWriter, r *http.Request) {
    // CORS — must check manually
    w.Header().Set("Access-Control-Allow-Origin", "*")

    // Auth — must check manually
    token := r.Header.Get("Authorization")
    user, err := validateJWT(token)
    if err != nil {
        http.Error(w, "unauthorized", 401)
        return
    }

    // Logging — must log manually
    log.Printf("GraphQL request from user %s", user.ID)

    // ... handle GraphQL query
}
// Every endpoint must redo all of this → easy to forget CORS or auth
```

```go
// WITH middleware — chain handles everything automatically
router.Use(
    corsMiddleware,       // Every request gets CORS headers
    authMiddleware,       // Every request is authenticated
    logMiddleware,        // Every request is logged
)
router.Handle("/query", graphqlHandler)  // Handler only needs to handle business logic
```

### Core Concept

```go
// Middleware chain
router := chi.NewRouter()
router.Use(
    corsMiddleware,           // 1. CORS headers
    authMiddleware,           // 2. JWT validation
    logMiddleware,            // 3. Request logging
    dataloaderMiddleware,     // 4. Initialize DataLoaders
    sentryMiddleware,         // 5. Error reporting
)
router.Handle("/query", graphqlHandler)
```

### When to Use
- HTTP-level concerns: CORS, auth, logging, tracing

### When NOT to Use
- GraphQL-specific concerns — use gql middleware (errors, recovery)

---

## Exercises

### EX1: Schema Analysis + Trace GraphQL Flow
1. Read `schema.graphqls`, list all Query and Mutation operations. For each:
   - Name, input types, return type, applied directives
2. Trace the flow of the `todos` query from schema through resolver -> usecase -> gRPC call -> todos service -> response
3. Draw a sequence diagram showing the complete flow

**Deliverable**: A markdown document with operation catalog and sequence diagram.

### EX2: Implement New Field with DataLoader
Add a new field `updatedAt: Time!` to the `Todo` type:
1. Update the GraphQL schema
2. Run `make generate`
3. Implement the resolver (if needed)
4. Update the mapper
5. Explain how DataLoader would work if this field required a separate service call. Analyze: how many gRPC calls would be made with and without DataLoader for a query listing 20 todos?

**Deliverable**: Working code changes (schema, resolver, mapper) + written DataLoader analysis.
