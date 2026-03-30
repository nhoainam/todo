# EX3 (Phase 2): Wire Dependency Graph

## Overview

This document traces the full Wire dependency graph starting from
`InitializeServer` in `di/wire.go`.

---

## WireSet Layout

```
di.InitializeServer
├── handler.WireSet
│   ├── handler.NewServer            → todov1.TodosServiceServer
│   │   ├── usecase.TodoGetter       (from service.WireSet)
│   │   ├── usecase.TodoUpdater      (from service.WireSet)
│   │   └── *validator.Validate      (from utils.WireSet)
│   └── grpc.NewServer               → (*grpc.Server, func(), error)
│       └── todov1.TodosServiceServer (from above)
│
├── service.WireSet
│   ├── service.NewTodoGetter        → usecase.TodoGetter
│   │   └── gateway.TodoQueriesGateway (from infra.WireSet)
│   ├── service.NewTodoUpdater       → usecase.TodoUpdater
│   │   ├── gateway.TodoQueriesGateway
│   │   └── gateway.TodoCommandsGateway (from infra.WireSet)
│   ├── service.NewTodoLister        → usecase.TodoLister
│   │   └── gateway.TodoQueriesGateway
│   ├── service.NewTodoCreator       → usecase.TodoCreator
│   │   ├── gateway.TodoCommandsGateway
│   │   └── idgen.IDGenerator        (from utils.WireSet)
│   └── service.NewTodoDeleter       → usecase.TodoDeleter
│       └── gateway.TodoCommandsGateway
│
├── infra.WireSet
│   ├── datastore.NewTodoReader      → gateway.TodoQueriesGateway
│   ├── datastore.NewTodoWriter      → gateway.TodoCommandsGateway
│   ├── datastore.NewBinder          → *datastore.Binder
│   │   └── *gorm.DB
│   └── gorm_app.Open                → (*gorm.DB, func(), error)
│       └── *config.Config
│
└── utils.WireSet
    ├── validator.New                → *validator.Validate
    └── idgen.NewIDGenerator         → idgen.IDGenerator
```

---

## Dependency Tree (simplified)

```
InitializeServer(*config.Config)
│
│  Wire-provided *config.Config
│
├─► gorm_app.Open(cfg) ──────────────────────────────► *gorm.DB
│
├─► datastore.NewBinder(db) ─────────────────────────► *datastore.Binder
│
├─► datastore.NewTodoReader() ───────────────────────► gateway.TodoQueriesGateway
│
├─► datastore.NewTodoWriter() ───────────────────────► gateway.TodoCommandsGateway
│
├─► idgen.NewIDGenerator() ──────────────────────────► idgen.IDGenerator
│
├─► service.NewTodoGetter(queriesGW) ────────────────► usecase.TodoGetter
├─► service.NewTodoUpdater(queriesGW, commandsGW) ───► usecase.TodoUpdater
├─► service.NewTodoLister(queriesGW) ────────────────► usecase.TodoLister
├─► service.NewTodoCreator(commandsGW, idgen) ───────► usecase.TodoCreator
├─► service.NewTodoDeleter(commandsGW) ──────────────► usecase.TodoDeleter
│
├─► validator.New() ─────────────────────────────────► *validator.Validate
│
├─► handler.NewServer(getter, updater, lister, creator, deleter, validate)
│                                                    ► todov1.TodosServiceServer
│
└─► grpc.NewServer(todosService) ────────────────────► *grpc.Server
```

---

## How Wire Resolves This

Wire uses **type matching**: each provider function declares its return type.
When another provider needs that type as a parameter, Wire automatically
connects them.

Example resolution chain for `usecase.TodoGetter`:
1. `service.NewTodoGetter` returns `usecase.TodoGetter` and needs
   `gateway.TodoQueriesGateway`.
2. `datastore.NewTodoReader` returns `gateway.TodoQueriesGateway` — Wire matches
   this automatically.
3. `datastore.NewTodoReader` needs no arguments (the DB comes from context) —
   the chain is complete.

### Compile-time guarantees
- Missing provider → **compile error** (not a runtime panic)
- Unused provider → **compile error** ("unused provider …")
- Circular dependency → **compile error**

---

## WireSet Locations

| Package | File | Provides |
|---|---|---|
| `di` | `wire.go` | `InitializeServer` injector |
| `internal/handler` | `wire.go` | `handler.WireSet` |
| `internal/service` | `wire.go` | `service.WireSet` |
| `internal/infra/datastore` | `wire.go` | `datastore.WireSet` |
| `internal/infra/gorm` | `wire.go` | `gorm.WireSet` |
| `internal/infra/idgen` | `wire.go` | `idgen.WireSet` |

---

## Running Wire

```bash
go tool wire ./di
```

This generates `di/wire_gen.go` with the full wiring code.
**Never edit `wire_gen.go` directly** — re-run `go tool wire` instead.
