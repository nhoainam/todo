// File: todos/internal/handler/grpc/grpc.go

package grpc

import (
	// Alias the standard grpc library to avoid a name collision with this
	// package, which is also named "grpc".
	grpclib "google.golang.org/grpc"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/handler/grpc/interceptor"
	todov1 "github.com/tuannguyenandpadcojp/fresher26/nam/todos/proto/todo/v1"
)

// NewServer builds a *grpc.Server wired with the Phase-1 interceptor chain.
//
// Parameters:
//   - todosService: the TodosService handler (internal/handler/todo_handler.go)
//   - dbInterceptor: injects *gorm.DB into every request context
//
// Returns:
//   - *grpclib.Server  — ready to call Serve(lis)
//   - func()           — cleanup; call on shutdown (flushes traces in later phases)
//   - error            — non-nil only if server construction itself fails
func NewServer(
	todosService todov1.TodosServiceServer,
	authInterceptor *interceptor.AuthInterceptor,
	dbInterceptor *interceptor.DBInterceptor,
) (*grpclib.Server, func(), error) {
	server := grpclib.NewServer(
		grpclib.ChainUnaryInterceptor(
			authInterceptor.Unary(), // Extract authenticated user id from metadata into context
			dbInterceptor.Unary(),   // Inject DB into every request context
		// grpctrace.UnaryServerInterceptor(...),     // 1. Datadog tracing
		//     logging.UnaryServerInterceptor(...),        // 2. Logging
		//     recovery.UnaryServerInterceptor(...),       // 3. Panic recovery
		//     sentryinterceptor.UnaryServerInterceptor(), // 4. Sentry error reporting
		//     authninterceptor.UnaryServerInterceptor(), // 5. Authentication
		//     authzinterceptor.UnaryServerInterceptor(), // 6. Authorization
		),
	)

	// Register the Todos service implementation.
	todov1.RegisterTodosServiceServer(server, todosService)

	cleanup := func() {
	}

	return server, cleanup, nil
}
