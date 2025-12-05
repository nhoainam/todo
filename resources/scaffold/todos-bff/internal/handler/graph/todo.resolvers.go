package graph

// todo.resolvers.go — Todo Query & Mutation Resolvers
//
// Phase 3: GraphQL & BFF Pattern
//
// This file is responsible for:
// 1. Implement query resolvers:
//    - Todo(ctx, name) — resolve a single todo by resource name
//    - Todos(ctx, listName, first, after) — resolve a paginated list
//
// 2. Implement mutation resolvers:
//    - CreateTodo(ctx, input) — create and return a new todo
//    - UpdateTodo(ctx, input) — update and return the todo
//    - DeleteTodo(ctx, name) — delete and return success payload
//
// 3. Implement field resolvers (lazy-loaded via DataLoader):
//    - Todo.Creator(ctx, obj) — resolve the creator User
//    - Todo.TodoList(ctx, obj) — resolve the parent TodoList
//
// Field resolvers are ONLY called if the client requests that field.
// Use DataLoader to batch multiple field resolver calls into one gRPC request.
//
// Pattern:
//   func (r *todoResolver) Creator(ctx context.Context, obj *model.Todo) (*model.User, error) {
//       return dataloader.For(ctx).UserLoader.Load(ctx, obj.CreatorID)()
//   }
//
// See: resources/phase-03-graphql-bff.md (resolvers, field resolvers, DataLoader)
