# Phase 1: Architecture & gRPC — Review Questions

## EX1: Trace Full Request Flow

1. When a GetTodo request arrives, which layers does it pass through in order? Why does the flow go in this specific direction?
2. What would happen if the handler layer directly imported the infrastructure layer (e.g., GORM)? What principle would that violate?
3. In the 5-step handler pattern (Parse → Build Input → Validate → Execute → Map Response), why is validation a separate step from parsing? What kind of errors does each step catch?
4. What is the purpose of the mapper between proto messages and domain entities? Why not pass proto messages directly to the use case?
5. How does an AppError (e.g., NewNotFound) get translated into a gRPC status code? What happens if you return a raw Go error instead of an AppError?

## EX2: Map Directory Structure and Identify Layers

1. If you needed to add a new entity (e.g., Tag), which directories would you create files in and why?
2. Why are gateway interfaces defined in a separate package from their implementations? What advantage does this give us?
3. Explain the difference between `internal/domain/` and `internal/infra/` in your own words. What belongs in each?
4. Given a random file in the project, how do you determine which Clean Architecture layer it belongs to? What clues do you look for?
5. Why does the domain layer have zero external dependencies? What problem does this solve?
6. What is "dependency inversion" and how is it applied in this project structure? Give a concrete example from the codebase.

## EX3: Local Dev Setup + List RPC Methods

1. Why do we use `envconfig` for configuration instead of hardcoding values or using a config file?
2. What is the purpose of the `Makefile` targets (build, test, lint, generate)? When would you use each one?
3. Why do we define the service contract in `.proto` files instead of directly in Go code? What problem does code generation solve?
4. What is the relationship between a proto `message` and a Go `struct`? Why don't we use the generated proto structs as our domain entities?
5. If you needed to add a new field to the Todo proto message, what steps would you take? What could break if you change an existing field number?

## EX4: Implement Mock Handler with 5-Step Pattern

1. Why does the handler depend on use case interfaces instead of concrete implementations? How does this help with testing?
2. Explain the resource name pattern `users/{user_id}/todo-lists/{list_id}/todos/{todo_id}`. Why use this format instead of separate ID parameters?
3. What happens when the use case returns an error? How should the handler translate different error types (NotFound, InvalidParameter, AuthZ) to gRPC responses?
4. Why are interceptors executed in a specific order (e.g., tracing → logging → panic recovery → auth)? What would go wrong if auth ran before panic recovery?
5. How does the authentication interceptor decide whether to allow or reject a request? What information does it need?
