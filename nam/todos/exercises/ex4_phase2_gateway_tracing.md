# EX1 (Phase 2): Gateway Tracing

## Overview

This document traces the `TodoQueriesGateway` interface through all layers,
mapping each interface method to its corresponding GORM query implementation.

---

## Interface Definition
**File**: `internal/domain/gateway/todo_queries.go`

```go
type TodoQueriesGateway interface {
    Get(ctx context.Context, todoID entity.TodoID, opts *GetTodoOptions) (*entity.Todo, error)
    List(ctx context.Context, opts *ListTodosOptions) (*OffsetPageResult[*entity.Todo], error)
}
```

---

## Implementation: `todoReader`
**File**: `internal/infra/datastore/todo_reader.go`

| Interface Method | GORM Gen Query |
|---|---|
| `Get(ctx, todoID, opts)` | `q.WithContext(ctx).Where(q.ID.Eq(int64(todoID))).First()` |
| `List(ctx, opts)` | `q.WithContext(ctx).Where(...filters...).Offset(offset).Limit(limit).FindByPage(offset, limit)` |

---

## Data Flow: Get

```
Handler (gRPC)
   │  req.Name → entity.ParseTodoResourceName → TodoResourceName.TodoID
   ▼
Service (todoGetter.Get)
   │  calls: todoQueriesGateway.Get(ctx, in.Name.TodoID, nil)
   ▼
Infrastructure (todoReader.Get)
   │  1. DBFromContext(ctx)        → *gorm.DB from context key
   │  2. query.Use(db).Todo        → typed query builder for todos table
   │  3. q.WithContext(ctx)        → adds context (cancellation, tracing)
   │  4. .Where(q.ID.Eq(todoID))  → WHERE id = ?   (parameterized, injection-safe)
   │  5. .First()                  → LIMIT 1, returns *model.Todo
   │  6. mapper.ToDomainTodo(m)    → converts model → entity
   ▼
Domain Entity
   └─ *entity.Todo returned up through service → handler → proto response
```

### GORM Gen SQL produced
```sql
SELECT `id`,`todo_list_id`,`title`,`content`,`status`,`due_date`,`priority`,
       `creator_id`,`created_at`,`updated_at`
FROM `todos`
WHERE `id` = ?
LIMIT 1
```

---

## Data Flow: List (with filter)

```
Handler (gRPC)
   │  builds ListTodosOptions{Filter: &TodoFilter{...}, Pagination: &OffsetPage{...}}
   ▼
Service (todoLister.List)
   │  calls: todoQueriesGateway.List(ctx, opts)
   ▼
Infrastructure (todoReader.List)
   │  1. DBFromContext(ctx)
   │  2. q.WithContext(ctx)
   │  3. applyFilter(q, qb, filter)
   │     ├─ filter.StatusEq    → .Where(q.Status.Eq(status.Int()))
   │     ├─ filter.PriorityEq  → .Where(q.Priority.Eq(priority.Int()))
   │     ├─ filter.DueDateGTE  → .Where(q.DueDate.Gte(*dueDateGTE))
   │     └─ filter.DueDateLTE  → .Where(q.DueDate.Lte(*dueDateLTE))
   │  4. .FindByPage(offset, limit) → returns ([]*model.Todo, count int64)
   │  5. maps each model.Todo → entity.Todo
   ▼
OffsetPageResult[*entity.Todo]{Items, TotalCount, Page}
```

### GORM Gen SQL produced (example with all filters)
```sql
SELECT `id`,`todo_list_id`,`title`,`content`,`status`,`due_date`,`priority`,
       `creator_id`,`created_at`,`updated_at`
FROM `todos`
WHERE `status` = ? AND `priority` = ? AND `due_date` >= ? AND `due_date` <= ?
LIMIT ? OFFSET ?
```

---

## Key Design Decisions

1. **Record-not-found is NOT an error**: `Get` returns `(nil, nil)` when no row
   matches — the service layer decides whether that means "NotFound error".
2. **DB lives in context**: Gateways call `DBFromContext(ctx)` so the same
   transaction is shared across multiple gateways in one use-case call.
3. **Model → Entity mapping is explicit**: A dedicated `mapper` package converts
   between the GORM model (`model.Todo`) and the domain entity (`entity.Todo`),
   keeping the domain layer free of GORM tags.
4. **Type-safe queries**: GORM Gen generates field helpers (`q.ID`, `q.Status`,
   etc.) so typos become compile errors, not runtime panics.
