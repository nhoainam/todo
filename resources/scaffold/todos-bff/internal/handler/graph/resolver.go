package graph

// resolver.go — Root Resolver
//
// Phase 3: GraphQL & BFF Pattern
//
// This file is responsible for:
// 1. Define the Resolver struct that holds all use case dependencies:
//    type Resolver struct {
//        TodoGetter  usecase.TodoGetter
//        TodoCreator usecase.TodoCreator
//        TodoUpdater usecase.TodoUpdater
//        TodoDeleter usecase.TodoDeleter
//        TodoLister  usecase.TodoLister
//    }
//
// 2. gqlgen uses this struct as the root — all resolvers receive it
//
// The root resolver is the "composition root" for GraphQL.
// Wire injects all dependencies into it.
//
// See: resources/phase-03-graphql-bff.md (resolver pattern)
