# Phase 4: Testing Patterns (Days 9-10)

## Goals
- Master unit testing with mocks (gomock)
- Understand integration testing patterns
- Master test infrastructure (template DB, bufconn)

---

## 1. Mock Generation with mockgen

### Problem Solved
Unit testing use cases requires isolation from databases and external APIs. You can't call a real database in unit tests because:
- Slow
- Requires database setup
- Non-deterministic (data changes)
- Can't test edge cases (network errors, timeouts)

**Real-world example:**
```go
// WITHOUT mocks — tests must call the real database
func TestTodoGetter_Get(t *testing.T) {
    // Must set up Postgres
    db, _ := gorm.Open(postgres.Open("host=localhost ..."))

    // Must insert test data
    db.Create(&Todo{ID: 1, Title: "Buy groceries"})

    getter := NewTodoGetter(NewTodoReader(db), ...)
    result, err := getter.Get(ctx, &input.TodoGetter{...})

    // Problem 1: Test takes 2-3 seconds (DB connection + query)
    // Problem 2: Run in parallel? DB conflicts!
    // Problem 3: Test "DB connection lost"? Must shut down Postgres!
    // Problem 4: CI/CD needs Postgres running → complex setup
}
```

```go
// WITH mocks — fast, isolated, deterministic tests
func TestTodoGetter_Get(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockGW := gatewaymock.NewMockTodoQueriesGateway(ctrl)
    mockGW.EXPECT().Get(gomock.Any(), entity.TodoID(1), gomock.Any()).
        Return(&entity.Todo{ID: 1, Title: "Buy groceries"}, nil)

    getter := NewTodoGetter(mockGW, ...)
    result, err := getter.Get(ctx, &input.TodoGetter{...})
    // Runs in < 1ms, no DB needed
    // Test error case? mockGW.EXPECT().Get(...).Return(nil, errors.New("connection lost"))
}
```

### Core Concept
`go.uber.org/mock` (mockgen) automatically generates mock implementations from interfaces.

#### Auto-generate mocks

```go
// File: todos/internal/domain/gateway/todo.go

//go:generate go tool mockgen -destination=mock/$GOFILE -source=$GOFILE

type TodoQueriesGateway interface {
    Get(ctx context.Context, todoID entity.TodoID, opts *GetTodoOptions) (*entity.Todo, error)
    List(ctx context.Context, opts *ListTodosOptions) (*query.OffsetPageResult[*entity.Todo], error)
}
```

Run `go generate ./...` -> creates file `domain/gateway/mock/todo.go`:
```go
// AUTO-GENERATED — DO NOT EDIT
type MockTodoQueriesGateway struct {
    ctrl     *gomock.Controller
    recorder *MockTodoQueriesGatewayMockRecorder
}

func NewMockTodoQueriesGateway(ctrl *gomock.Controller) *MockTodoQueriesGateway { ... }
```

#### Using mocks in tests

```go
func TestTodoGetter_Get(t *testing.T) {
    // 1. Create mock controller
    ctrl := gomock.NewController(t)

    // 2. Create mocks
    mockTodoGateway := gatewaymock.NewMockTodoQueriesGateway(ctrl)
    mockClock := timemock.NewMockClock(ctrl)

    // 3. Set up expectations
    mockTodoGateway.EXPECT().
        Get(gomock.Any(), entity.TodoID(456), gomock.Any()).
        Return(expectedTodo, nil)

    mockClock.EXPECT().
        Now().
        Return(fixedTime)

    // 4. Create SUT (System Under Test) with mocks
    getter := NewTodoGetter(mockTodoGateway, ..., mockClock, ...)

    // 5. Execute
    result, err := getter.Get(ctx, &input.TodoGetter{...})

    // 6. Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)

    // gomock automatically verifies all expectations when the test ends
}
```

### Mock Matchers

```go
// Exact match
mockGW.EXPECT().Get(gomock.Any(), entity.TodoID(456), gomock.Any())

// Any value
mockGW.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any())

// Custom matcher
mockGW.EXPECT().BatchGet(gomock.Any(), gomock.InAnyOrder([]entity.TodoID{1, 2, 3}))

// Multiple calls
mockGW.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Times(3)

// Sequential returns
mockGW.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
    Return(todo1, nil).
    Return(todo2, nil)

// Return error
mockGW.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
    Return(nil, errors.New("db connection lost"))
```

### When to Use
- Unit testing usecases/services
- Testing error paths (DB errors, not found, permission denied)
- Testing business logic in isolation

### When NOT to Use
- Integration tests (use real dependencies)
- Testing the infrastructure layer (GORM queries) — use a real DB
- Mocking too much — if you need to mock 10+ dependencies, the architecture may have issues

---

## 2. Table-Driven Tests

### Problem Solved
Writing a separate test function for each case:
- Duplicated code
- Hard to see how many cases there are
- Hard to add new cases

**Real-world example:**
```go
// WITHOUT table-driven — each case is a separate function
func TestGetTodo_Success(t *testing.T) {
    ctrl := gomock.NewController(t)                    // duplicate
    mockGW := gatewaymock.NewMockTodoQueriesGateway(ctrl)  // duplicate
    mockGW.EXPECT().Get(gomock.Any(), entity.TodoID(1), gomock.Any()).Return(todo, nil)
    getter := NewTodoGetter(mockGW, ...)               // duplicate
    result, err := getter.Get(ctx, &input.TodoGetter{Name: name1})
    assert.NoError(t, err)                             // duplicate
    assert.Equal(t, expected, result)
}

func TestGetTodo_NotFound(t *testing.T) {
    ctrl := gomock.NewController(t)                    // AGAIN duplicate!
    mockGW := gatewaymock.NewMockTodoQueriesGateway(ctrl)  // AGAIN duplicate!
    mockGW.EXPECT().Get(gomock.Any(), entity.TodoID(999), gomock.Any()).Return(nil, nil)
    getter := NewTodoGetter(mockGW, ...)               // AGAIN duplicate!
    _, err := getter.Get(ctx, &input.TodoGetter{Name: name999})
    assert.Error(t, err)
}

// Adding a "permission denied" case? Copy-paste all setup code again...
// 10 cases → 10 functions with 80% identical code
```

Table-driven tests consolidate all cases in one place, differing only in input/expected/mock setup.

### Core Concept
Group all test cases into a map/slice and loop through them.

```go
func Test_batchTodosGetter_BatchGet(t *testing.T) {
    // Shared test data
    var (
        now        = time.Now()
        userID     = shared.UserID(100)
        todoID1    = todos.TodoID(456)
        todoName1  = todos.NewTodoResourceName(userID, todoListID, todoID1)

        entityTodo1 = &entity.Todo{
            ID:         todoID1.ToEntity(),
            TodoListID: cast.Ptr(todoListID.ToEntity()),
            Title:      "Buy groceries",
            Status:     entity.TodoStatusPending,
        }
    )

    // Mock field types
    type fields struct {
        mockTodoGateway   *gatewaymock.MockTodoQueriesGateway
        mockPermChecker   *portmock.MockTodoPermissionChecker
        mockClock         *timemock.MockClock
    }

    tests := map[string]struct {
        prepare  func(f *fields)    // Set up mock expectations
        args     *input.BatchTodosGetter
        expected *output.BatchTodosGetter
        wantErr  bool
    }{
        "success: get todo": {
            prepare: func(f *fields) {
                f.mockPermChecker.EXPECT().
                    CanView(gomock.Any(), todoName1).
                    Return(entityTodo1, nil)
            },
            args:     &input.BatchTodosGetter{Names: []todos.TodoResourceName{todoName1}},
            expected: &output.BatchTodosGetter{Todos: []*todos.Todo{modelTodo1}},
        },
        "error: permission denied": {
            prepare: func(f *fields) {
                f.mockPermChecker.EXPECT().
                    CanView(gomock.Any(), todoName1).
                    Return(nil, errors.NewAuthZ("no access", nil))
            },
            args:    &input.BatchTodosGetter{Names: []todos.TodoResourceName{todoName1}},
            wantErr: true,
        },
        "success: todo not found returns partial results": {
            prepare: func(f *fields) {
                f.mockPermChecker.EXPECT().
                    CanView(gomock.Any(), todoName1).
                    Return(nil, errors.NewNotFound("not found", nil))
            },
            args:     &input.BatchTodosGetter{Names: []todos.TodoResourceName{todoName1}},
            expected: &output.BatchTodosGetter{Todos: []*todos.Todo{}},
        },
    }

    for name, tt := range tests {
        t.Run(name, func(t *testing.T) {
            t.Parallel()

            // Setup
            ctrl := gomock.NewController(t)
            f := &fields{
                mockTodoGateway:  gatewaymock.NewMockTodoQueriesGateway(ctrl),
                mockPermChecker: portmock.NewMockTodoPermissionChecker(ctrl),
                mockClock:       timemock.NewMockClock(ctrl),
            }
            tt.prepare(f)

            // Create SUT
            sut := NewBatchTodosGetter(f.mockTodoGateway, f.mockPermChecker, f.mockClock, cfg)

            // Execute
            got, err := sut.BatchGet(ctx, tt.args)

            // Assert
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            if diff := cmp.Diff(tt.expected, got); diff != "" {
                t.Errorf("mismatch (-want +got):\n%s", diff)
            }
        })
    }
}
```

### When to Use
- Every test function with multiple cases (> 2)
- Cases differ in input/expected but share the same flow

### When NOT to Use
- Only 1-2 simple cases
- Cases with completely different flows

---

## 3. Deterministic Mocks (Time, ID)

### Problem Solved
Tests using `time.Now()` or random ID generation give different results each run.

**Real-world example:**
```go
// WITHOUT time mock — non-deterministic test
func (s *service) CreateTodo(ctx context.Context, in *input.Creator) (*output.Creator, error) {
    todo := &entity.Todo{
        Title:     in.Title,
        CreatedAt: time.Now(),  // Different value each run!
    }
    return s.gateway.Create(ctx, todo)
}

// Test:
func TestCreateTodo(t *testing.T) {
    result, _ := service.CreateTodo(ctx, input)
    assert.Equal(t, expected.CreatedAt, result.CreatedAt)
    // FAIL! expected: 2025-01-01T10:00:00 but got: 2025-01-01T10:00:01
    // Because time.Now() runs at 2 different moments

    // IDs are similar:
    // expected.ID = 12345 but got.ID = 67890 (random snowflake ID)
}
```

Mock Clock and IDGenerator ensure tests always produce the same results.

### Core Concept

#### Time Mock
```go
// Interface
type Clock interface {
    Now() time.Time
}

// Production implementation
type realClock struct{}
func (c *realClock) Now() time.Time { return time.Now() }

// Test: mock clock
mockClock := timemock.NewMockClock(ctrl)
fixedTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
mockClock.EXPECT().Now().Return(fixedTime).AnyTimes()
```

#### ID Mock
```go
// Interface
type IDGenerator interface {
    Generate() (int64, error)
}

// Test: mock ID generator
mockIDGen := idgenmock.NewMockIDGenerator(ctrl)
mockIDGen.EXPECT().Generate().Return(int64(12345), nil)
```

### When to Use
- Whenever code uses `time.Now()` or random values
- Assertions need deterministic values

### When NOT to Use
- Tests that don't care about the specific value of time/ID

---

## 4. Integration Tests

### Problem Solved
Unit tests with mocks don't guarantee:
- GORM queries actually run correctly
- SQL schema matches Go structs
- Transaction behavior is correct
- Full request/response flow works

**Real-world example:**
```go
// Unit test passes — but production fails!

// Mock test: gateway.Get() returns todo → test passes ✓
mockGW.EXPECT().Get(gomock.Any(), id, gomock.Any()).Return(todo, nil)
// But in reality: the GORM query has a bug!

// Actual bug in gateway implementation:
func (r *todoReader) Get(ctx context.Context, id entity.TodoID) (*entity.Todo, error) {
    return query.Use(db).Todo.WithContext(ctx).
        Where(query.Todo.TodoListID.Eq(int64(id))).  // BUG! Uses TodoListID instead of ID
        First()
}
// Unit test DOESN'T catch this bug because the mock bypasses the entire GORM layer
// Integration test with a real DB will catch it: "expected todo ID=456, got todo with TodoListID=456"

// Similarly: wrong struct tag "gorm:column:titl" (typo) → unit test passes, production fails
```

### Core Concept
Integration tests run with a real database and in-memory gRPC server.

#### Test Lifecycle

```go
func TestMain(m *testing.M) {
    // 1. Create template database (once for all tests)
    tmplDB, err := testutil.InitTemplateDB(context.TODO())
    if err != nil { log.Fatalf("cannot init template db: %s", err) }

    // 2. Run all tests
    exitVal := m.Run()

    // 3. Cleanup
    tmplDB.Release(context.TODO())
    os.Exit(exitVal)
}
```

#### Template Database Pattern

```
TestMain:
  ┌─ Create template DB
  ├─ Apply schema (from SQL file)
  └─ Load fixtures (from YAML)

Test1: Clone template → Run test → Cleanup clone
Test2: Clone template → Run test → Cleanup clone
Test3: Clone template → Run test → Cleanup clone
```

**Why clone?** Each test gets its own database → parallel testing, no interference.

#### Test Helper

```go
type TodosServiceTestHelper struct {
    client todospb.TodosServiceClient  // gRPC client
    ctx    context.Context             // Context with auth
    gormDB *gorm.DB                   // Direct DB access for setup/verification
}

func NewTodosServiceTestHelper(
    t *testing.T,
    currentUser shared.User,
    mockCtrl *gomock.Controller,
    mockServices *MockServices,
) *TodosServiceTestHelper {
    // 1. Create test DB (clone from template)
    gormDB := testutil.InitDB(t)

    // 2. Set up config
    cfg := &config.Config{ServerPort: 50051, ...}

    // 3. Initialize server via Wire (with mocks for external services)
    grpcServer, cleanup, _ := registry.InitializeServer(cfg, logger)
    t.Cleanup(cleanup)

    // 4. Create in-memory gRPC connection (bufconn)
    lis := bufconn.Listen(bufSize)
    go grpcServer.Serve(lis)

    conn, _ := grpc.Dial("bufnet",
        grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
            return lis.DialContext(ctx)
        }),
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )

    // 5. Create client with auth metadata
    client := todospb.NewTodosServiceClient(conn)
    md := metadata.Pairs("x-authenticated-user-id", currentUser.Name.UserID.String())
    ctx := metadata.NewOutgoingContext(context.Background(), md)

    return &TodosServiceTestHelper{client: client, ctx: ctx, gormDB: gormDB}
}
```

#### Integration Test Example

```go
func TestGetTodo(t *testing.T) {
    t.Parallel()

    tests := map[string]struct {
        prepare  func(t *testing.T, h *TodosServiceTestHelper, ms *MockServices)
        args     *todospb.GetTodoRequest
        expected *todospb.GetTodoResponse
        wantErr  codes.Code
    }{
        "success": {
            prepare: func(t *testing.T, h *TodosServiceTestHelper, ms *MockServices) {
                // Insert test data into the real DB
                todo := &entity.Todo{ID: 1, Title: "Buy groceries", Status: 0}
                h.gormDB.Create(todo)
            },
            args: &todospb.GetTodoRequest{
                Name: "users/1/todo-lists/1/todos/1",
            },
            expected: &todospb.GetTodoResponse{
                Todo: &todospb.Todo{Name: "users/1/todo-lists/1/todos/1", ...},
            },
        },
        "not found": {
            prepare: func(t *testing.T, h *TodosServiceTestHelper, ms *MockServices) {
                // Don't insert data -> todo doesn't exist
            },
            args: &todospb.GetTodoRequest{
                Name: "users/1/todo-lists/1/todos/999",
            },
            wantErr: codes.NotFound,
        },
    }

    for name, tt := range tests {
        t.Run(name, func(t *testing.T) {
            t.Parallel()
            ctrl := gomock.NewController(t)
            ms := NewMockServices(ctrl)
            h := NewTodosServiceTestHelper(t, currentUser, ctrl, ms)

            tt.prepare(t, h, ms)

            got, err := h.client.GetTodo(h.ctx, tt.args)

            if tt.wantErr != 0 {
                assert.Equal(t, tt.wantErr, status.Code(err))
                return
            }
            assert.NoError(t, err)
            if diff := cmp.Diff(tt.expected, got, protocmp.Transform()); diff != "" {
                t.Errorf("mismatch (-want +got):\n%s", diff)
            }
        })
    }
}
```

### When to Use
- Testing full gRPC request/response flow
- Testing database queries (GORM)
- Testing authorization logic end-to-end
- Testing data integrity (foreign keys, unique constraints)

### When NOT to Use
- Testing simple business logic (unit tests are faster)
- Testing external APIs (mock them in integration tests)

---

## 5. bufconn — In-Memory gRPC

### Problem Solved
Integration tests need a gRPC server but you don't want to:
- Bind a real TCP port (conflicts with other tests)
- Have network overhead
- Deal with port cleanup issues

**Real-world example:**
```go
// Real TCP port — many problems
func TestGetTodo(t *testing.T) {
    lis, _ := net.Listen("tcp", ":50051")  // Bind port 50051
    go grpcServer.Serve(lis)

    conn, _ := grpc.Dial("localhost:50051", ...)
    client := todospb.NewTodosServiceClient(conn)
    // ... test
}

// Problem 1: Run 2 tests simultaneously? "port 50051 already in use"!
// Problem 2: Test fails midway → port isn't freed → next test fails
// Problem 3: CI runs multiple test suites in parallel → port conflicts
// Problem 4: Network latency (even localhost) → tests slower than necessary
```

```go
// bufconn — in-memory, no port needed
lis := bufconn.Listen(1024 * 1024)  // In-memory buffer
go grpcServer.Serve(lis)
// No port binding → no conflicts
// In-memory → faster than TCP
// Automatically cleaned up when the test ends
```

### Core Concept
`bufconn` creates an in-memory connection that behaves like TCP but doesn't need a port.

```go
lis := bufconn.Listen(1024 * 1024)  // 1MB buffer

// Server side
go grpcServer.Serve(lis)

// Client side
conn, _ := grpc.Dial("bufnet",
    grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
        return lis.DialContext(ctx)
    }),
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
client := todospb.NewTodosServiceClient(conn)
```

### When to Use
- Every gRPC integration test
- Testing the interceptor chain

### When NOT to Use
- Testing network-specific behavior (timeouts, reconnection)

---

## 6. Assertion Patterns

### Problem Solved
Comparing complex structs with `assert.Equal` produces hard-to-read output when tests fail.

**Real-world example:**
```go
// assert.Equal — hard-to-read output with large structs
assert.Equal(t, expected, got)
// Output on failure:
// Expected: &{ID:1 Title:Buy groceries Status:0 Priority:1 TodoListID:5 CreatedAt:2025-01-01 00:00:00 +0000 UTC
//   UpdatedAt:2025-01-01 00:00:00 +0000 UTC DueDate:2025-03-15 AssigneeID:789 ...}
// Actual:   &{ID:1 Title:Buy vegetables Status:0 Priority:1 TodoListID:5 CreatedAt:2025-01-01 00:00:00 +0000 UTC
//   UpdatedAt:2025-01-01 00:00:00 +0000 UTC DueDate:2025-03-15 AssigneeID:789 ...}
// Two 5-line structs — where's the difference? Must read each field to find it!
// Only "Buy groceries" vs "Buy vegetables" but the output shows the ENTIRE struct
```

```go
// cmp.Diff — only shows the differing fields
if diff := cmp.Diff(expected, got); diff != "" {
    t.Errorf("mismatch (-want +got):\n%s", diff)
}
// Output:
//   Todo{
// -   Title: "Buy groceries",
// +   Title: "Buy vegetables",
//   }
// Clear! Only shows the different field, doesn't spam identical fields
```

### Core Concept

#### cmp.Diff — Readable diffs
```go
if diff := cmp.Diff(expected, got); diff != "" {
    t.Errorf("mismatch (-want +got):\n%s", diff)
}
```

#### protocmp.Transform — Proto comparison
```go
// Proto messages may have unexported fields, need protocmp
if diff := cmp.Diff(expected, got, protocmp.Transform()); diff != "" {
    t.Errorf("mismatch (-want +got):\n%s", diff)
}
```

#### testify assertions
```go
assert.NoError(t, err)
assert.Error(t, err)
assert.Equal(t, expected, got)
assert.Nil(t, result)
assert.NotNil(t, result)
assert.Contains(t, slice, element)
```

### When to Use
- `cmp.Diff` for complex struct comparison
- `protocmp.Transform()` for protobuf messages
- `testify/assert` for simple assertions

---

## Exercises

### EX1: Write Unit Test with Table-Driven Pattern
Write a unit test for a use case (e.g., `TodoGetter.Get`) using table-driven test pattern:
- Happy path: todo exists
- Error: todo not found
- Error: permission denied
- Error: gateway returns error

Use gomock for all dependencies. Follow the table-driven pattern with `map[string]struct{}` and `t.Run()`.

**Deliverable**: Working test file with at least 4 test cases using table-driven pattern.

### EX2: Write Integration Test with bufconn + Template DB
1. Read existing integration test files to understand the setup
2. Add a new test case to an existing integration test (e.g., add "get todo with tags" case to `get_todo_test`)
3. The test must use:
   - Template database pattern (clone DB per test)
   - bufconn for in-memory gRPC connection
   - Real database queries (no mocks for datastore layer)

**Deliverable**: Working integration test case + written explanation of how template DB and bufconn work together.

### EX3: Mock Exercise with Deterministic Patterns
1. Create a mock for a new interface:
   - Add `//go:generate` directive
   - Run `go generate`
2. Write a test using the mock with deterministic time and ID:
   - Mock `Clock` to return a fixed time
   - Mock `IDGenerator` to return a fixed ID
   - Verify the output uses the mocked values

**Deliverable**: Working mock generation + test file demonstrating deterministic patterns.
