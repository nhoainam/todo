package service

import "github.com/google/wire"

// WireSet provides all service (use-case) implementations.
var WireSet = wire.NewSet(
	NewTodoGetter,
	NewTodoUpdater,
	NewTodoLister,
	NewTodoCreator,
	NewTodoDeleter,
)
