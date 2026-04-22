package grpc

import (
	"github.com/google/wire"
	"github.com/tuannguyenandpadcojp/fresher26/nam/users/internal/handler/grpc/interceptor"
)

var WireSet = wire.NewSet(NewServer, interceptor.NewDBInterceptor)
