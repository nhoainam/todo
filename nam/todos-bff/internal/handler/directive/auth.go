package directive

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/apperrors"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/handler/graph/model"
)

// Phase 3: @hasPermission directive implementation (authorization check on GraphQL fields). See resources/phase-03-graphql-bff.md
// File: todos-bff/internal/handler/graph/directives/

func hasRequiredPermissions(ctx context.Context, required []model.Permission) bool {
	// Placeholder: Extract user permissions from context (e.g., from JWT claims)
	userPermissions := []model.Permission{"TODO_DELETE"} // Example permissions
	for _, rp := range required {
		found := false
		for _, up := range userPermissions {
			if up == rp {
				found = true
				break
			}
		}
		if !found {
			return false // Missing required permission
		}
	}
	return true // All required permissions are present
}

func HasPermission() func(ctx context.Context, obj any, next graphql.Resolver, permissions []model.Permission) (any, error) {
	return func(ctx context.Context, obj any, next graphql.Resolver, permissions []model.Permission) (any, error) {
		// Check user permissions
		if !hasRequiredPermissions(ctx, permissions) {
			return nil, apperrors.NewAuthZ("insufficient permissions", nil)
		}
		return next(ctx) // Continue to resolver
	}
}

// func ValidateInput(v *validator.Validate) func(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
// 	return func(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
// 		// Validate input using go-playground/validator
// 		return next(ctx)
// 	}
// }
