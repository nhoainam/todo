# Fresher26 — Go Training: Todos Project

A hands-on training program for Go freshers to learn the production tech stack through building a Todos application (backend + BFF).

**Duration**: 2.5 weeks (13 working days) across 5 phases

## Prerequisites

- Basic Go knowledge (syntax, goroutines, interfaces)
- Go 1.21+ installed
- Docker & Docker Compose
- `make` available in your shell

## Getting Started

### 1. Clone the repository

```bash
git clone <repo-url>
cd fresher26
```

### 2. Scaffold your workspace

Each fresher has their own directory. Run:

```bash
make scaffold FRESHER=<your-name>
```

This copies the project scaffold into `<your-name>/todos/` (backend) and `<your-name>/todos-bff/` (BFF). Existing files are never overwritten.

### 3. Follow the training plan

Open `resources/training-plan.md` for the full schedule, coding conventions, and graduation checklist.

## Training Phases

| Phase | Days | Topic | Exercises | Materials |
|-------|------|-------|-----------|-----------|
| **1** | 1-3 | Architecture & gRPC | 4 | `resources/phase-01-architecture-grpc.md` |
| **2** | 4-6 | Database & DI | 4 | `resources/phase-02-database-di.md` |
| **3** | 7-8 | GraphQL & BFF | 2 | `resources/phase-03-graphql-bff.md` |
| **4** | 9-10 | Testing Patterns | 3 | `resources/phase-04-testing.md` |
| **5** | 11-13 | Observability & Capstone | 1 | `resources/phase-05-observability-e2e.md` |

## Directory Structure

```
fresher26/
├── resources/
│   ├── training-plan.md          # Full training plan & conventions
│   ├── phase-01-architecture-grpc.md
│   ├── phase-02-database-di.md
│   ├── phase-03-graphql-bff.md
│   ├── phase-04-testing.md
│   ├── phase-05-observability-e2e.md
│   ├── review-questions/         # Questions for PR reviews
│   └── scaffold/                 # Project template (todos + todos-bff)
├── <your-name>/
│   ├── todos/                    # Your backend service (gRPC)
│   ├── todos-bff/                # Your BFF service (GraphQL)
│   └── progress.md               # Your progress tracker
└── .github/
    └── pull_request_template.md  # PR template — fill this out for every PR
```

## Submitting PRs

1. Create a branch from `main`: `feat/<exercise>-short-description`
2. Implement the exercise in your directory
3. Fill out the PR template — include your phase, exercise number, and answers to review questions from `resources/review-questions/phase-{N}.md`
4. Request a review

## Tech Stack

| Layer | Backend (todos) | BFF (todos-bff) |
|-------|-----------------|------------------|
| Transport | gRPC + Protobuf | GraphQL (gqlgen) + HTTP (chi) |
| Architecture | Clean Architecture | Clean Architecture |
| DI | Google Wire | Google Wire |
| Database | GORM + MySQL | Calls backend via gRPC |
| Testing | testify + gomock | testify + gomock |

## Resources

- **Training plan**: `resources/training-plan.md`
- **Phase materials**: `resources/phase-{N}-*.md`
- **Review questions**: `resources/review-questions/phase-{N}.md`
- **Scaffold source**: `resources/scaffold/`
