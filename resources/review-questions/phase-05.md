# Phase 5: Observability & End-to-End — Review Questions

## EX: Implement "Add Tag to Todo" End-to-End

1. When implementing a new feature end-to-end, what is the order you create files in and why? (e.g., do you start from domain or handler?) Explain your reasoning.
2. How does structured logging (Zap) differ from `fmt.Println` or `log.Println`? Why is structured logging important in production?
3. What is distributed tracing and how does Datadog connect a request across the BFF and backend services? What identifier links them?
4. When should an error be reported to Sentry vs. just logged? Give examples of each.
5. After implementing the feature, what code generation commands did you run and in what order? Why does the order matter?
6. How does your new feature integrate with the existing Wire dependency graph? What WireSets did you modify?
