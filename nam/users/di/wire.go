//go:build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/config"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/handler"
	grpcserver "github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/handler/grpc"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/infra/datastore"
	gorm_app "github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/infra/gorm"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/service"
	grpclib "google.golang.org/grpc"
)

func InitializeServer(cfg *config.Config) (*grpclib.Server, func(), error) {
	wire.Build(
		// Infrastructure
		datastore.WireSet, // NewUserReader (datastore.UserReader)
		gorm_app.WireSet,  // Open (*gorm.DB)

		// Services (use-case implementations)
		service.WireSet, // NewUserService (*service.UserService)

		// Handler
		handler.WireSet,

		// gRPC server
		grpcserver.WireSet, // NewServer (*grpc.Server)
	)
	return nil, nil, nil
}
