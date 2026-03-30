package idgen

import "github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"

// IDGenerator generates new unique TodoIDs.
// Keeping it as an interface allows swapping implementations (UUID, snowflake, ULID, etc.)
// without touching the use-case layer.
type IDGenerator interface {
	NewTodoID() entity.TodoID
}
