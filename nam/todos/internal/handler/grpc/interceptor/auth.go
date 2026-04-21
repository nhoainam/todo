package interceptor

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const AuthenticatedUserIDHeader = "x-authenticated-user-id"

type authenticatedUserIDKey struct{}

type AuthInterceptor struct{}

func NewAuthInterceptor() *AuthInterceptor {
	return &AuthInterceptor{}
}

func WithAuthenticatedUserID(ctx context.Context, userID entity.UserID) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if userID <= 0 {
		return ctx
	}

	return context.WithValue(ctx, authenticatedUserIDKey{}, userID)
}

func AuthenticatedUserIDFromContext(ctx context.Context) (entity.UserID, bool) {
	if ctx == nil {
		return 0, false
	}

	userID, ok := ctx.Value(authenticatedUserIDKey{}).(entity.UserID)
	if !ok || userID <= 0 {
		return 0, false
	}

	return userID, true
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		userID, err := authenticatedUserIDFromMetadata(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		ctx = WithAuthenticatedUserID(ctx, userID)
		return handler(ctx, req)
	}
}

func authenticatedUserIDFromMetadata(ctx context.Context) (entity.UserID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, fmt.Errorf("missing incoming metadata")
	}

	values := md.Get(AuthenticatedUserIDHeader)
	if len(values) == 0 {
		return 0, fmt.Errorf("%s header is required", AuthenticatedUserIDHeader)
	}

	raw := strings.TrimSpace(values[0])
	if raw == "" {
		return 0, fmt.Errorf("%s header must not be empty", AuthenticatedUserIDHeader)
	}

	parsed, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s header: %w", AuthenticatedUserIDHeader, err)
	}
	if parsed <= 0 {
		return 0, fmt.Errorf("invalid %s header: must be > 0", AuthenticatedUserIDHeader)
	}

	return entity.UserID(parsed), nil
}
