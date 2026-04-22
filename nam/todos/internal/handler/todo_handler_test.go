package handler_test

// todo_handler_test.go — Unit Tests for GetTodo gRPC Handler
//
// ─────────────────────────────────────────────────────────────────────────────
// INTERCEPTOR CHAIN — execution order before GetTodo handler is reached
// ─────────────────────────────────────────────────────────────────────────────
//
// When a client calls GetTodo, the request travels through this chain BEFORE the
// handler function body runs (see internal/handler/grpc/grpc.go):
//
//   Client Request
//       │
//       ▼
//   [1] Tracing Interceptor (grpctrace.UnaryServerInterceptor)
//       – Creates a distributed trace span for the RPC call.
//       – Injects the trace/span IDs into the context.
//       – Currently commented-out in Phase 1; added in Phase 5 (observability).
//       │
//       ▼
//   [2] Logging Interceptor (logging.UnaryServerInterceptor)
//       – Logs the incoming RPC method name, request metadata, and latency.
//       – Runs *around* the handler so it can log both start and completion.
//       – Currently commented-out; added in Phase 5.
//       │
//       ▼
//   [3] Recovery Interceptor (recovery.UnaryServerInterceptor)
//       – Catches any panic that occurs in downstream interceptors or the handler.
//       – Converts the panic into a gRPC codes.Internal error instead of crashing
//         the server process.
//       – Currently commented-out; added in Phase 5.
//       │
//       ▼
//   [4] Sentry Interceptor (sentryinterceptor.UnaryServerInterceptor)
//       – Reports unexpected errors to Sentry for alerting and aggregation.
//       – Runs after recovery so it sees the recovered error, not the raw panic.
//       – Currently commented-out; added in Phase 5.
//       │
//       ▼
//   [5] Authentication Interceptor (authninterceptor.UnaryServerInterceptor)
//       – Reads the "x-authenticated-user-id" header from gRPC metadata.
//       – Loads the full user object and injects it into the context.
//       – Returns codes.Unauthenticated if the header is missing or invalid.
//       – Currently commented-out; added in Phase 2/3.
//       │
//       ▼
//   [6] Authorization Interceptor (authzinterceptor.UnaryServerInterceptor)
//       – Reads the authenticated user from context (set by step 5).
//       – Checks whether the user has permission to call this specific RPC.
//       – Returns codes.PermissionDenied if access is not allowed.
//       – Currently commented-out; added in Phase 2/3.
//       │
//       ▼
//   [7] handler.GetTodo  ← THIS IS WHERE OUR TESTS EXERCISE LOGIC
//       – Parse → Build Input → Validate → Execute → Map Response
//
// Because unit tests exercise the handler directly (bypassing the gRPC server),
// interceptors [1]–[6] do NOT run in these tests. This is intentional: unit
// tests focus on handler logic only. Integration/e2e tests cover the full chain.
//
// ─────────────────────────────────────────────────────────────────────────────

import (
	"context"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/apperrors"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/handler"
	grpcinterceptor "github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/handler/grpc/interceptor"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/input"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/usecase/output"
	todov1 "github.com/tuannguyenandpadcojp/fresher26/nam/todos/proto/todo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ─────────────────────────────────────────────────────────────────────────────
// Mock implementations
// ─────────────────────────────────────────────────────────────────────────────

// mockTodoGetter is a hand-written mock of the usecase.TodoGetter interface.
// It captures the last call's input and returns the pre-configured response or error.
type mockTodoGetter struct {
	calledWith *input.TodoGetter
	returnTodo *entity.Todo
	returnErr  error
}

func (m *mockTodoGetter) Get(_ context.Context, in *input.TodoGetter) (*output.TodoGetter, error) {
	m.calledWith = in
	if m.returnErr != nil {
		return nil, m.returnErr
	}
	return &output.TodoGetter{Todo: m.returnTodo}, nil
}

// mockTodoUpdater satisfies the usecase.TodoUpdater interface with a no-op.
type mockTodoUpdater struct{}

func (m *mockTodoUpdater) Update(_ context.Context, in *input.TodoUpdater) (*output.TodoUpdater, error) {
	return nil, nil
}

type mockTodoLister struct {
	calledWith *input.TodoLister
	returnOut  *output.TodoLister
	returnErr  error
}

func (m *mockTodoLister) List(_ context.Context, in *input.TodoLister) (*output.TodoLister, error) {
	m.calledWith = in
	if m.returnErr != nil {
		return nil, m.returnErr
	}
	if m.returnOut != nil {
		return m.returnOut, nil
	}
	return &output.TodoLister{}, nil
}

// mockTodoCreator satisfies usecase.TodoCreator with a no-op.
type mockTodoCreator struct{}

func (m *mockTodoCreator) Create(_ context.Context, _ *input.TodoCreator) (*output.TodoCreator, error) {
	return nil, nil
}

// mockTodoDeleter satisfies usecase.TodoDeleter with a no-op.
type mockTodoDeleter struct{}

func (m *mockTodoDeleter) Delete(_ context.Context, _ *input.TodoDeleter) error {
	return nil
}

// Compile-time interface checks — fails to compile if the mocks drift from the real interfaces.
var _ usecase.TodoGetter = (*mockTodoGetter)(nil)
var _ usecase.TodoUpdater = (*mockTodoUpdater)(nil)
var _ usecase.TodoLister = (*mockTodoLister)(nil)
var _ usecase.TodoCreator = (*mockTodoCreator)(nil)
var _ usecase.TodoDeleter = (*mockTodoDeleter)(nil)

// ─────────────────────────────────────────────────────────────────────────────
// Test helpers
// ─────────────────────────────────────────────────────────────────────────────

// newTestServer constructs a handler.server ready for unit testing.
func newTestServer(getter usecase.TodoGetter) todov1.TodosServiceServer {
	return handler.NewServer(getter, &mockTodoUpdater{}, &mockTodoLister{}, &mockTodoCreator{}, &mockTodoDeleter{}, validator.New())
}

func authenticatedContext(userID entity.UserID) context.Context {
	return grpcinterceptor.WithAuthenticatedUserID(context.Background(), userID)
}

// mustStatusCode asserts the gRPC status code of an error and returns it.
func mustStatusCode(t *testing.T, err error, want codes.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %v, got nil", want)
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got: %v", err)
	}
	if st.Code() != want {
		t.Fatalf("expected gRPC code %v, got %v (message: %q)", want, st.Code(), st.Message())
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Test cases
// ─────────────────────────────────────────────────────────────────────────────

// TestGetTodo_HappyPath validates the full successful flow:
//
//	Step 1 Parse      — valid resource name is parsed into TodoResourceName
//	Step 2 BuildInput — TodoGetter input DTO is constructed
//	Step 3 Validate   — struct validation passes
//	Step 4 Execute    — usecase returns a todo entity
//	Step 5 Map        — entity is mapped to a proto response
func TestGetTodo_HappyPath(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	mockGetter := &mockTodoGetter{
		returnTodo: &entity.Todo{
			ID:        entity.TodoID(456),
			ListID:    entity.TodoListID(200),
			CreatorID: entity.UserID(100),
			Title:     "Buy groceries",
			Content:   "Milk, eggs, bread",
			Status:    entity.TodoStatusPENDING,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	srv := newTestServer(mockGetter)

	resp, err := srv.GetTodo(authenticatedContext(100), &todov1.GetTodoRequest{
		Name: "users/100/todo-lists/200/todos/456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify Step 2: input DTO was built with the correct resource name.
	if mockGetter.calledWith == nil {
		t.Fatal("usecase.Get was not called")
	}
	if mockGetter.calledWith.Name.TodoID != 456 {
		t.Errorf("TodoID: want 456, got %d", mockGetter.calledWith.Name.TodoID)
	}
	if mockGetter.calledWith.Name.UserID != 100 {
		t.Errorf("UserID: want 100, got %d", mockGetter.calledWith.Name.UserID)
	}
	if mockGetter.calledWith.Name.TodoListID != 200 {
		t.Errorf("TodoListID: want 200, got %d", mockGetter.calledWith.Name.TodoListID)
	}

	// Verify Step 5: proto response is mapped correctly.
	if resp.GetTodo() == nil {
		t.Fatal("response Todo is nil")
	}
	if resp.GetTodo().Title != "Buy groceries" {
		t.Errorf("Title: want %q, got %q", "Buy groceries", resp.GetTodo().Title)
	}
	if resp.GetTodo().Status != todov1.TodoStatus_PENDING {
		t.Errorf("Status: want PENDING, got %v", resp.GetTodo().Status)
	}
	wantName := "users/100/todo-lists/200/todos/456"
	if resp.GetTodo().Name != wantName {
		t.Errorf("Name: want %q, got %q", wantName, resp.GetTodo().Name)
	}
}

// TestGetTodo_InvalidResourceName validates Step 1 (Parse) error handling.
// When the resource name is malformed the handler must return codes.InvalidArgument
// without calling the usecase at all.
func TestGetTodo_InvalidResourceName(t *testing.T) {
	cases := []struct {
		name        string
		requestName string
	}{
		{"empty name", ""},
		{"wrong format", "todos/456"},
		{"non-numeric id", "users/abc/todo-lists/200/todos/456"},
		{"zero user id", "users/0/todo-lists/200/todos/456"},
		{"zero todo list id", "users/100/todo-lists/0/todos/456"},
		{"zero todo id", "users/100/todo-lists/200/todos/0"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockGetter := &mockTodoGetter{}
			srv := newTestServer(mockGetter)

			_, err := srv.GetTodo(context.Background(), &todov1.GetTodoRequest{
				Name: tc.requestName,
			})

			mustStatusCode(t, err, codes.InvalidArgument)

			// Usecase must NOT have been called — the error was caught in Step 1.
			if mockGetter.calledWith != nil {
				t.Error("usecase.Get should not have been called for an invalid resource name")
			}
		})
	}
}

func TestGetTodo_Unauthenticated(t *testing.T) {
	mockGetter := &mockTodoGetter{}
	srv := newTestServer(mockGetter)

	_, err := srv.GetTodo(context.Background(), &todov1.GetTodoRequest{
		Name: "users/1/todo-lists/1/todos/1",
	})

	mustStatusCode(t, err, codes.Unauthenticated)
	if mockGetter.calledWith != nil {
		t.Error("usecase.Get should not have been called when request is unauthenticated")
	}
}

func TestGetTodo_AuthenticatedUserMismatch(t *testing.T) {
	mockGetter := &mockTodoGetter{}
	srv := newTestServer(mockGetter)

	_, err := srv.GetTodo(authenticatedContext(999), &todov1.GetTodoRequest{
		Name: "users/1/todo-lists/1/todos/1",
	})

	mustStatusCode(t, err, codes.PermissionDenied)
	if mockGetter.calledWith != nil {
		t.Error("usecase.Get should not have been called for forbidden resource access")
	}
}

// TestGetTodo_NotFound validates Step 4 (Execute) error mapping.
// When the usecase returns an apperrors.NotFound error the handler must
// translate it to gRPC codes.NotFound.
func TestGetTodo_NotFound(t *testing.T) {
	mockGetter := &mockTodoGetter{
		returnErr: apperrors.NewNotFound("todo not found", nil),
	}

	srv := newTestServer(mockGetter)

	_, err := srv.GetTodo(authenticatedContext(1), &todov1.GetTodoRequest{
		Name: "users/1/todo-lists/1/todos/999",
	})

	mustStatusCode(t, err, codes.NotFound)
}

// TestGetTodo_InternalError validates Step 4 error mapping for unexpected errors.
// Any non-AppError (e.g., a raw database error) must become codes.Internal so
// that implementation details are never leaked to the caller.
func TestGetTodo_InternalError(t *testing.T) {
	mockGetter := &mockTodoGetter{
		returnErr: apperrors.NewInternal("database connection lost", nil),
	}

	srv := newTestServer(mockGetter)

	_, err := srv.GetTodo(authenticatedContext(1), &todov1.GetTodoRequest{
		Name: "users/1/todo-lists/1/todos/1",
	})

	mustStatusCode(t, err, codes.Internal)
}

// TestGetTodo_AuthZError validates that a permission-denied error from the
// usecase is correctly surfaced as codes.PermissionDenied.
func TestGetTodo_AuthZError(t *testing.T) {
	mockGetter := &mockTodoGetter{
		returnErr: apperrors.NewAuthZ("user does not have access to this todo", nil),
	}

	srv := newTestServer(mockGetter)

	_, err := srv.GetTodo(authenticatedContext(1), &todov1.GetTodoRequest{
		Name: "users/1/todo-lists/1/todos/1",
	})

	mustStatusCode(t, err, codes.PermissionDenied)
}

func TestListTodos_ParentListNotFound(t *testing.T) {
	listUsecase := &mockTodoLister{returnErr: apperrors.NewNotFound("todo list not found", nil)}
	srv := handler.NewServer(&mockTodoGetter{}, &mockTodoUpdater{}, listUsecase, &mockTodoCreator{}, &mockTodoDeleter{}, validator.New())

	_, err := srv.ListTodos(authenticatedContext(5), &todov1.ListTodosRequest{
		Parent: "users/5/todo-lists/2",
		Limit:  20,
		Offset: 0,
	})

	mustStatusCode(t, err, codes.NotFound)
	if listUsecase.calledWith == nil {
		t.Fatal("usecase.List should be called")
	}
}
