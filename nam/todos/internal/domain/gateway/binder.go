package gateway

import "context"

// Binder handles database transaction scoping within context boundaries.
// It is implemented in the infrastructure layer.
type Binder interface {
	// Bind starts a new transaction and embeds it into the returned context.
	Bind(ctx context.Context) context.Context
	
	// Commit commits the transaction associated with the context.
	Commit(ctx context.Context) error
	
	// Rollback rolls back the transaction associated with the context.
	Rollback(ctx context.Context) error
}
