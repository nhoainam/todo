# Phase 3: GraphQL & BFF — Review Questions

## EX1: Schema Analysis + Trace GraphQL Flow

1. What is the BFF pattern and why do we need a separate service between the frontend and backend? Why not have the frontend call gRPC directly?
2. How does GraphQL differ from REST in terms of data fetching? What problem does GraphQL solve that REST doesn't?
3. What is the purpose of custom scalars like `ResourceName` and `Time`? Why not use plain `String` and `Int`?
4. Trace a GraphQL query from schema → resolver → use case → gRPC client. How does this compare to the backend's request flow?
5. What is a field resolver and when does it get called? What happens if a client doesn't request a field that has a resolver?
6. How does the BFF translate gRPC errors from the backend into GraphQL errors for the frontend?

## EX2: Implement New Field with DataLoader

1. What is the N+1 problem? Give a concrete example with todos and their creators.
2. How does DataLoader solve the N+1 problem? Explain the batching mechanism and the timing window.
3. Why must DataLoaders be created per-request (not shared across requests)? What would happen if they were shared?
4. After adding a new field to the GraphQL schema, what commands do you run and what files get generated? What do you implement manually vs. what is auto-generated?
5. What is a GraphQL directive (e.g., `@hasPermission`)? How is it different from middleware?
6. If you add a field resolver that makes a gRPC call, how would you optimize it to avoid N+1? What changes would you need?
