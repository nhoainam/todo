package idgen

import "github.com/google/wire"

// WireSet provides the IDGenerator to the Wire dependency graph.
var WireSet = wire.NewSet(NewIDGenerator)
