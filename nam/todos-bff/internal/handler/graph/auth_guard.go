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
