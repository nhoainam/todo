package grpc

import (
	"github.com/google/wire"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/handler/grpc/interceptor"
)

// WireSet provides the gRPC server constructor and its interceptors.
var WireSet = wire.NewSet(NewServer, interceptor.NewDBInterceptor)
