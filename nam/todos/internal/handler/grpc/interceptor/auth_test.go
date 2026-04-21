package interceptor

import (
	"context"
	"testing"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuthInterceptorUnary_SetsAuthenticatedUserInContext(t *testing.T) {
	unary := NewAuthInterceptor().Unary()
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(AuthenticatedUserIDHeader, "123"))

	handlerCalled := false
	resp, err := unary(ctx, "req", &grpc.UnaryServerInfo{}, func(ctx context.Context, req any) (any, error) {
		handlerCalled = true

		gotUserID, ok := AuthenticatedUserIDFromContext(ctx)
		if !ok {
			t.Fatal("expected authenticated user id in context")
		}
		if gotUserID != entity.UserID(123) {
			t.Fatalf("expected authenticated user id 123, got %d", gotUserID)
		}
		if req != "req" {
			t.Fatalf("expected req to be forwarded to handler, got %#v", req)
		}

		return "ok", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !handlerCalled {
		t.Fatal("expected handler to be called")
	}
	if resp != "ok" {
		t.Fatalf("expected response %q, got %#v", "ok", resp)
	}
}

func TestAuthInterceptorUnary_MissingHeaderReturnsUnauthenticated(t *testing.T) {
	unary := NewAuthInterceptor().Unary()
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-some-header", "value"))

	handlerCalled := false
	_, err := unary(ctx, nil, &grpc.UnaryServerInfo{}, func(ctx context.Context, req any) (any, error) {
		handlerCalled = true
		return nil, nil
	})
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("expected code %v, got %v (err=%v)", codes.Unauthenticated, status.Code(err), err)
	}
	if handlerCalled {
		t.Fatal("handler should not be called when authentication header is missing")
	}
}

func TestAuthInterceptorUnary_InvalidHeaderReturnsUnauthenticated(t *testing.T) {
	unary := NewAuthInterceptor().Unary()
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(AuthenticatedUserIDHeader, "abc"))

	handlerCalled := false
	_, err := unary(ctx, nil, &grpc.UnaryServerInfo{}, func(ctx context.Context, req any) (any, error) {
		handlerCalled = true
		return nil, nil
	})
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("expected code %v, got %v (err=%v)", codes.Unauthenticated, status.Code(err), err)
	}
	if handlerCalled {
		t.Fatal("handler should not be called when authentication header is invalid")
	}
}
