package datastore

import "github.com/google/wire"

// WireSet provides the datastore gateway implementations.
var WireSet = wire.NewSet(
	NewTodoReader,
	NewTodoWriter,
	NewBinder,
)
