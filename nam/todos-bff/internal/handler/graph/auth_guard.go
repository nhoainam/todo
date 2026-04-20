package graph

import (
	"context"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/apperrors"
	"github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/domain/entity"
	http_middleware "github.com/tuannguyenandpadcojp/fresher26/nam/todos-bff/internal/middleware/http"
)

func requireAuthenticatedUserID(ctx context.Context) (entity.UserID, error) {
	userID, ok := http_middleware.UserIDFromContext(ctx)
	if !ok {
		return 0, apperrors.NewAuthN("authentication required", nil)
	}

	return entity.UserID(userID), nil
}

func ensureResourceUserMatch(ctx context.Context, resourceUserID entity.UserID) error {
	if resourceUserID <= 0 {
		return apperrors.NewInvalidParameter("invalid resource user id", nil)
	}

	authenticatedUserID, err := requireAuthenticatedUserID(ctx)
	if err != nil {
		return err
	}

	if authenticatedUserID != resourceUserID {
		return apperrors.NewAuthZ(
			"forbidden resource access",
			nil,
			apperrors.WithMetadata("resource_user_id", resourceUserID.Int64()),
			apperrors.WithMetadata("context_user_id", authenticatedUserID.Int64()),
		)
	}

	return nil
}
