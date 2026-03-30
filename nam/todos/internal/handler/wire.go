package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

// NewValidator creates a new input validator instance.
func NewValidator() *validator.Validate {
	return validator.New()
}

// WireSet provides the gRPC handler (TodosServiceServer).
var WireSet = wire.NewSet(
	NewServer,
	NewValidator,
)
