package idgen

import (
	"math/rand"

	"github.com/tuannguyenandpadcojp/fresher26/nam/todos/internal/domain/entity"
)

type randomIDGenerator struct{}

// NewIDGenerator returns an IDGenerator backed by a random int64.
// In production you would use a distributed ID generator (e.g. Snowflake / ULID).
func NewIDGenerator() IDGenerator {
	return &randomIDGenerator{}
}

func (g *randomIDGenerator) NewTodoID() entity.TodoID {
	return entity.TodoID(rand.Int63()) //nolint:gosec // non-security ID generation
}
