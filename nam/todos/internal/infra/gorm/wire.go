package gorm_app

import "github.com/google/wire"

// WireSet provides the *gorm.DB instance.
var WireSet = wire.NewSet(Open)
