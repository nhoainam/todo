package service

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewLoginService,
	NewUserLogoutService,
	NewUserRegisterService,
)
