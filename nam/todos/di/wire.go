//go:build wireinject

package di

import (
	"github.com/google/wire"
	grpclib "google.golang.org/grpc"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/config"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/handler"
	grpcserver "github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/handler/grpc"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/infra/datastore"
	gorm_app "github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/infra/gorm"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/infra/idgen"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/service"
)

// Phase 2: Wire dependency injection (providers, WireSets, injector). See resources/phase-02-database-di.md

// InitializeServer builds and returns a fully wired *grpc.Server.
// Wire generates the implementation of this function in wire_gen.go.
func InitializeServer(cfg *config.Config) (*grpclib.Server, func(), error) {
	wire.Build(
		// Infrastructure
		gorm_app.WireSet,  // Open(*gorm.DB)
		datastore.WireSet, // NewTodoReader, NewTodoWriter, NewBinder
		idgen.WireSet,     // NewIDGenerator (EX4: new provider)

		// Services (use-case implementations)
		service.WireSet, // NewTodoGetter, NewTodoUpdater, NewTodoLister, NewTodoCreator, NewTodoDeleter

		// Handler
		handler.WireSet, // NewServer (TodosServiceServer)

		// gRPC server
		grpcserver.WireSet, // NewServer (*grpc.Server)
	)
	return nil, nil, nil
}
