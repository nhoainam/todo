## Week & Exercise

<!-- Which week and exercise is this PR for? Example: Week 2 — BT3: Interceptor Analysis -->

**Week**:
**Exercise**:

## What I Implemented

<!-- Bullet list of what you built or changed -->

-

## Explain Your Decisions

<!-- Pick 2-3 questions from the question bank for your exercise and answer them.
     Question bank: resources/review-questions/week-XX.md (replace XX with your week number) -->

### Question 1:
<!-- Your answer -->

### Question 2:
<!-- Your answer -->

### Question 3 (optional):
<!-- Your answer -->

## Self-Check

<!-- Check each item that applies to your PR. Leave unchecked if not applicable to this exercise. -->

- [ ] Files are in the correct Clean Architecture layer (no cross-layer imports)
- [ ] Handler follows the 5-step pattern (Parse → Build Input → Validate → Execute → Map Response)
- [ ] Gateway interfaces separate Commands (write) and Queries (read)
- [ ] Unit tests use table-driven pattern with `prepare`, `args`, `expected`, `wantErr`
- [ ] IDs use strong types (`TodoID`, `UserID`, `TodoListID`), not raw strings
- [ ] Error handling uses `AppError` (`NewNotFound`, `NewInvalidParameter`, etc.)
- [ ] Wire DI setup is correct (WireSet defined, providers registered)

## Notes (optional)

<!-- Anything else: what you're unsure about, what you'd like feedback on, etc. -->
