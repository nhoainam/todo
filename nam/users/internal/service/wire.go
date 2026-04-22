package service

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewLoginService,
	NewTokenGeneratorService,
	NewUserLogoutService,
	NewUserRegisterService,
)
