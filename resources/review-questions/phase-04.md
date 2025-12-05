# Phase 4: Testing Patterns — Review Questions

## EX1: Write Unit Test with Table-Driven Pattern

1. Why do we use table-driven tests instead of writing separate test functions for each case? What are the benefits?
2. Explain the `prepare` function in a table-driven test. Why do we set up mock expectations there instead of inline?
3. What is the difference between `gomock.Any()` and setting a specific expected value? When should you use each?

## EX2: Write Integration Test with bufconn + Template DB

1. What is the template database pattern? Why create a template and clone it per-test instead of using one shared database?
2. What is `bufconn` and why do we use it for integration tests? What advantage does it have over starting a real gRPC server on a TCP port?
3. How does the test helper struct (`TodosServiceTestHelper`) simplify writing integration tests? What would the tests look like without it?
4. When adding a new test case to an existing table-driven test, how do you decide what mock expectations to set? How do you know which methods will be called?
5. Why do we use `cmp.Diff` instead of `assert.Equal` for comparing complex structs? What does the output look like when there's a difference?

## EX3: Mock Exercise with Deterministic Patterns

1. How does `mockgen` know which methods to generate mocks for? What does the `//go:generate` directive do?
2. What is the difference between `gomock.InOrder()` and `gomock.InAnyOrder()`? When does the order of mock calls matter?
3. Why do we mock time (`Clock`) and ID generation (`IDGenerator`) in tests? What would happen if we didn't?
