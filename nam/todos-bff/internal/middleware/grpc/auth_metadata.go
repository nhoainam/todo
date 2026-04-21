package grpc_middleware

import (
	"context"
	"strconv"

	http_middleware "github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/middleware/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const AuthenticatedUserIDHeader = "x-authenticated-user-id"

func NewAuthMetadataUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		userID, ok := http_middleware.UserIDFromContext(ctx)
		if ok {
			ctx = metadata.AppendToOutgoingContext(ctx, AuthenticatedUserIDHeader, strconv.FormatInt(userID, 10))
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
