# Go Fresher Training Plan — Todos Project

> **Target audience**: Developers who know basic Go (syntax, goroutines, interfaces), new to the project's tech stack.
> **Duration**: 2.5 weeks (13 working days), organized into 5 phases
> **Goal**: Be able to independently implement a new feature end-to-end on both backend (gRPC) and BFF (GraphQL gateway).
> **Reference**: Patterns and conventions in this training are referenced from real services in the monorepo (e.g., `go/services/knowledge`).

---

## Tech Stack Overview

| Layer | Backend Service | BFF |
|---|---|---|
| **Transport** | gRPC + Protobuf | GraphQL (gqlgen) + HTTP (chi) |
| **Architecture** | Clean Architecture (Domain -> UseCase -> Handler -> Infrastructure) | Same Clean Architecture |
| **DI** | Google Wire | Google Wire |
| **Database** | GORM + MySQL | None directly (calls via gRPC) |
| **Auth** | gRPC interceptors (authn/authz) | HTTP middleware (JWT/Auth0) |
| **Observability** | Datadog tracing, Sentry, Zap logging | Datadog, Sentry, Zap |
| **Validation** | go-playground/validator | go-playground/validator |
| **Testing** | testify + gomock + integration tests | testify + gomock |
| **Code Gen** | protoc, wire, mockgen, enumer | gqlgen, wire, mockgen |
| **Config** | envconfig | envconfig |
| **ID Generation** | Snowflake/ULID | -- |
| **Deployment** | Docker + Helm + K8s | Docker + Helm + K8s |

---

## Project Structure

### Backend Service — `todos/`

```
todos/
├── cmd/
│   └── todos/main.go              # gRPC server entry point
├── internal/
│   ├── config/                    # App configuration (envconfig)
│   ├── domain/
│   │   ├── entity/                # Domain models (Todo, TodoList, Tag, User...)
│   │   ├── gateway/               # Repository interfaces (Commands + Queries)
│   │   ├── model/                 # Value objects, resource names
│   │   └── query/                 # Pagination, sorting helpers
│   ├── handler/
│   │   ├── grpc/                  # gRPC handlers + interceptors
│   │   │   ├── interceptor/       # authn, authz, logging, recovery, sentry
│   │   │   ├── mapper/            # Proto <-> Domain mapping
│   │   │   └── service/           # gRPC service implementations
│   ├── infrastructure/
│   │   ├── datastore/             # GORM repository implementations
│   │   └── ...                    # Other external service clients
│   ├── service/                   # Use case implementations + helpers
│   ├── usecase/todos/             # Use case interfaces
│   │   ├── input/                 # Use case input DTOs
│   │   └── output/                # Use case output DTOs
│   ├── registry/                  # Wire DI setup
│   └── utils/                     # Shared utilities
├── database/
│   └── schemas/                   # DB schema definitions (Ridgepole)
└── test/
    └── integration/               # Integration tests (gRPC)
```

### BFF — `todos-bff/`

```
todos-bff/
├── cmd/todos-bff/main.go          # HTTP + GraphQL server entry point
├── graph/
│   └── schema.graphqls            # GraphQL schema definition
├── internal/
│   ├── config/                    # App configuration
│   ├── domain/
│   │   ├── gateway/               # Gateway interfaces (gRPC backend calls)
│   │   ├── model/                 # Domain models
│   │   └── query/                 # Query helpers
│   ├── handler/graph/             # GraphQL resolvers
│   │   ├── generated/             # gqlgen generated code
│   │   ├── model/                 # GraphQL models
│   │   ├── dataloaders/           # N+1 prevention
│   │   ├── directives/            # Custom GraphQL directives
│   │   ├── mapper/                # GraphQL <-> Domain mapping
│   │   └── scalar/                # Custom scalars (Date, ResourceName)
│   ├── infrastructure/            # gRPC client implementations
│   ├── middleware/
│   │   ├── gql/                   # GraphQL middleware (errors, logging, recovery)
│   │   └── http/                  # HTTP middleware (auth, cors, dataloaders)
│   ├── service/                   # Service layer
│   ├── usecase/                   # Use case layer (input/output/mock)
│   └── registry/                  # Wire DI setup
└── test/                          # Tests
```

---

## Training Schedule

| Phase | Days | Topic | Exercises | Detailed Materials |
|---|---|---|---|---|
| **1** | Days 1-3 | Architecture & gRPC | 4 | [phase-01-architecture-grpc.md](./phase-01-architecture-grpc.md) |
| **2** | Days 4-6 | Database & DI | 4 | [phase-02-database-di.md](./phase-02-database-di.md) |
| **3** | Days 7-8 | GraphQL & BFF | 2 | [phase-03-graphql-bff.md](./phase-03-graphql-bff.md) |
| **4** | Days 9-10 | Testing Patterns | 3 | [phase-04-testing.md](./phase-04-testing.md) |
| **5** | Days 11-13 | Observability & Capstone | 1 capstone | [phase-05-observability-e2e.md](./phase-05-observability-e2e.md) |

### Coach Checkpoints

| Checkpoint | When | What to Review |
|---|---|---|
| **Checkpoint 1** | Day 3 (end of Phase 1) | Fresher can trace request flows, understands 4 layers, local dev works. Review EX1-EX4 deliverables. |
| **Checkpoint 2** | Day 6 (end of Phase 2) | Fresher can write GORM queries, understands Wire DI, gateway pattern is clear. Review EX1-EX4 deliverables. |
| **Checkpoint 3** | Day 10 (end of Phase 4) | Fresher can write unit + integration tests, understands GraphQL/BFF. Review all Phase 3-4 deliverables. |

> **If a fresher falls behind**: extend by 1-2 days at the checkpoint boundary rather than skipping content. The coach determines readiness to advance at each checkpoint.

---

## Coding Conventions

### Naming
- Gateway interfaces: `{Entity}CommandsGateway`, `{Entity}QueriesGateway`
- UseCase interfaces: `{Action}{Entity}` (e.g., `TodoGetter`, `BatchTodosCreator`)
- Input/Output DTOs: `input.{UseCaseName}`, `output.{UseCaseName}`
- Mock files: `mock/{filename}.go` (auto-generated)
- Enum generation: `go:generate enumer --type=...`

### Error Handling
- Custom error types in `internal/errors/`
- gRPC status codes mapping (InvalidArgument, NotFound, PermissionDenied...)
- Wrap errors with context: `fmt.Errorf("get todo: %w", err)`

### Validation
- Use `go-playground/validator` tags on input structs
- Validate in the handler layer before calling use cases

### Wire
- Each package has a `WireSet` variable
- Parent packages aggregate child WireSets
- `cleanup` function for resource cleanup (DB connections, etc.)

---

## Graduation Checklist

After completing all 5 phases, a fresher should be able to:

- [ ] Explain Clean Architecture and the dependency flow in the project
- [ ] Implement a new gRPC endpoint end-to-end (proto -> handler -> usecase -> gateway -> test)
- [ ] Implement a new GraphQL resolver in the BFF (schema -> resolver -> gRPC call -> test)
- [ ] Write unit tests with mocks and integration tests
- [ ] Understand Wire DI and add new dependencies to the graph
- [ ] Read and debug production issues via logs/traces
- [ ] Run the code generation workflow (`make generate`)
- [ ] Review teammates' PRs based on project conventions
