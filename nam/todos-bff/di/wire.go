//go:build wireinject

package di

import (
	"context"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/google/wire"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/handler/dataloader"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/handler/directive"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/handler/graph"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/handler/graph/generated"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/infra/grpc_client"
	grpc_middleware "github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/middleware/grpc"
	http_middleware "github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/middleware/http"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/service"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/usecase"
	todopb "github.com/tuannguyenandpadcojp/fresher26/nam/todos/proto/todo/v1"
	userpb "github.com/tuannguyenandpadcojp/fresher26/nam/users/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

// Phase 2: Wire dependency injection for BFF (gRPC clients, resolvers, dataloaders). See resources/phase-02-database-di.md

var appSet = wire.NewSet(
	NewTodosGRPCConn,
	NewTodosServiceClient,
	grpc_client.NewTodoServiceClient,
	service.NewTodoService,
	service.NewTodoUpdater,
	service.NewAuthLoginService,
	service.NewAuthLogoutService,
	service.NewAuthRegisterService,
	http_middleware.NewJWTTokenVerifier,
	wire.Bind(new(http_middleware.TokenVerifier), new(*http_middleware.JWTTokenVerifier)),
	NewUserGRPCConn,
	NewUserServiceClient,
	grpc_client.NewAuthServiceClient,
	NewResolver,
	NewExecutableSchema,
	NewGraphQLServer,
	NewHTTPHandler,
)

// InitializeServer builds and returns a fully wired GraphQL HTTP handler.
func InitializeServer() (http.Handler, func(), error) {
	wire.Build(appSet)
	return nil, nil, nil
}

type UserGRPCConn struct {
	Conn *grpc.ClientConn
}

type TodosGRPCConn struct {
	Conn *grpc.ClientConn
}

func NewUserGRPCConn() (*UserGRPCConn, func(), error) {
	target := os.Getenv("USERS_GRPC_ADDR")
	if target == "" {
		target = os.Getenv("USERS_SERVICE_ADDR")
	}
	if target == "" {
		target = "localhost:50052"
	}

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = conn.Close()
	}

	return &UserGRPCConn{Conn: conn}, cleanup, nil
}

func NewUserServiceClient(conn *UserGRPCConn) userpb.UserServiceClient {
	return userpb.NewUserServiceClient(conn.Conn)
}

func NewTodosGRPCConn() (*TodosGRPCConn, func(), error) {
	target := os.Getenv("TODOS_GRPC_ADDR")
	if target == "" {
		target = os.Getenv("TODOS_SERVICE_ADDR")
	}
	if target == "" {
		target = "localhost:50051"
	}

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(grpc_middleware.NewAuthMetadataUnaryClientInterceptor()),
	)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = conn.Close()
	}

	return &TodosGRPCConn{Conn: conn}, cleanup, nil
}

func NewTodosServiceClient(conn *TodosGRPCConn) todopb.TodosServiceClient {
	return todopb.NewTodosServiceClient(conn.Conn)
}

func NewResolver(
	todoGetter usecase.TodoGetter,
	todoUpdater usecase.TodoUpdater,
	authLogin usecase.AuthLogin,
	authLogout usecase.AuthLogout,
	authRegister usecase.AuthRegister,
) *graph.Resolver {
	return &graph.Resolver{
		TodoGetter:   todoGetter,
		TodoUpdater:  todoUpdater,
		AuthLogin:    authLogin,
		AuthLogout:   authLogout,
		AuthRegister: authRegister,
	}
}

func NewExecutableSchema(resolver *graph.Resolver) graphql.ExecutableSchema {
	cfg := generated.Config{Resolvers: resolver}
	cfg.Directives.HasPermission = directive.HasPermission()
	cfg.Directives.ValidateInput = passthroughValidationDirective

	return generated.NewExecutableSchema(cfg)
}

func passthroughValidationDirective(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
	return next(ctx)
}

func NewGraphQLServer(schema graphql.ExecutableSchema) *gqlhandler.Server {
	return gqlhandler.NewDefaultServer(schema)
}

func NewHTTPHandler(server *gqlhandler.Server, tokenVerifier http_middleware.TokenVerifier) http.Handler {
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := dataloader.WithLoaders(r.Context(), dataloader.NewLoaders())
		server.ServeHTTP(w, r.WithContext(ctx))
	})

	withHTTPContext := http_middleware.InjectHTTPMiddleware()(baseHandler)
	return http_middleware.AuthMiddleware(tokenVerifier)(withHTTPContext)
}
