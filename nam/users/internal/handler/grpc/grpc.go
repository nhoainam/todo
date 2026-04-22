package grpc

import (
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/handler/grpc/interceptor"
	userv1 "github.com/tuannguyenandpadcojp/fresher26/nam/users/proto/user/v1"
	grpclib "google.golang.org/grpc"
)

func NewServer(userService userv1.UserServiceServer, dbInterceptor *interceptor.DBInterceptor) (*grpclib.Server, func(), error) {
	server := grpclib.NewServer(
		grpclib.ChainUnaryInterceptor(
			dbInterceptor.Unary(), // Inject DB into every request context
		),
	)

	userv1.RegisterUserServiceServer(server, userService)

	cleanup := func() {
	}

	return server, cleanup, nil
}
