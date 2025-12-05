# Phase 2: Database & DI — Review Questions

## EX1: Gateway Tracing

1. Why do we separate TodoCommandsGateway (write) and TodoQueriesGateway (read) instead of having one TodoGateway? What benefit does this give us?
2. Where is the gateway interface defined vs. where is it implemented? Why this separation?
3. What is the "database context pattern" (DB stored in context)? Why pass the DB connection through context instead of storing it in the struct?

## EX2: Write GORM Query

1. Why do we use a separate GORM model (in `persistence/model/`) instead of using the domain entity directly with GORM? What problems does this prevent?
2. What is the purpose of the mapper between GORM models and domain entities? When does this conversion happen?
3. How does GORM's soft delete (`DeletedAt`) work? What query does GORM generate differently when soft delete is enabled?

## EX3: Wire Dependency Graph

1. What is the difference between Wire's compile-time DI and runtime DI (like `dig`)? Why does this project use Wire?
2. What is a WireSet and why does each package define one? How do they compose together?
3. If Wire fails to compile, what are the most common causes? How do you debug a Wire error?

## EX4: Add New Provider

1. When you add a new provider to a WireSet, what other changes are needed for Wire to include it in the dependency graph?
2. What happens if two providers return the same interface type? How does Wire resolve ambiguity?
3. Explain the relationship between the `wire.go` file (with build tag `//go:build wireinject`) and the generated `wire_gen.go`. Why do we need both?
