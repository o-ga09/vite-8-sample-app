# Project Overview

A multi-tenant household budget management system built with Go (backend) and React/Vite (frontend).

| Layer | Technology |
|-------|------------|
| ORM | [bob](https://github.com/stephenafamo/bob) v0.42 |
| Database | PostgreSQL 17 |
| API | [ogen](https://github.com/ogen-go/ogen) (OpenAPI v3 → Go server) |
| Migration | [dbmate](https://github.com/amacneil/dbmate) |
| Backend Lint | [golangci-lint](https://golangci-lint.run/) |
| Frontend | React 19 + Vite 8 + TypeScript strict |
| UI | Tailwind CSS v4 + shadcn/ui |
| API Client | openapi-fetch (generated from `api/openapi.yaml`) |
| Logging | slog (JSON) → Promtail → Loki → Grafana |
| Toolchain | [mise](https://mise.jdx.dev/) (Go, Node.js, pnpm) |

## Key Directories

```
api/openapi.yaml          # OpenAPI schema (shared between backend and frontend)
cmd/server/main.go        # API server entry point (default port: 18080)
cmd/seed/main.go          # Sample data seeder
db/migrations/            # dbmate migrations
internal/
  handler/                # ogen handler implementations (DTO conversion)
  infra/db/               # DB connection, WorkspaceScopedExec, WithTx
  infra/dbgen/            # bob-generated code (models, where, joins, etc.)
  infra/hook/             # QueryHooks (multi-tenant auto filter)
  oas/                    # ogen-generated code
  repository/             # Data access layer (bob query building)
  service/                # Business logic (transaction management)
queries/                  # Hand-written SQL + bob-generated code
webapp/                   # React/Vite frontend
```

---

# Development Workflow

## Initial Setup

Prerequisites: [mise](https://mise.jdx.dev/) and Docker must be installed.

```sh
mise trust && mise install          # Install Go, Node.js, pnpm via mise
docker compose up -d                # Start PostgreSQL, Loki, Grafana, Promtail
mise run migrate                    # Run database migrations
mise run seed                       # Insert sample data
```

## Backend Commands

All backend commands run from the **repository root**.

```sh
# Run
mise run server                     # Start API server (port 18080)
go run ./cmd/server/                # Alternative: run directly with go

# Test
go test ./...                       # Run all Go tests
go test ./internal/...              # Run tests in a specific package

# Lint
mise run lint                       # golangci-lint run
go tool golangci-lint run           # Alternative: run directly

# Code generation (requires DB connection)
mise run bobgen                     # Regenerate bob ORM code from DB schema
mise run ogen                       # Regenerate ogen server code from OpenAPI spec
```

### Port Configuration

Default ports (override in `mise.local.toml`):

| Service | Port |
|---------|------|
| API server | 18080 |
| PostgreSQL | 15432 |
| Grafana | 13000 |
| Loki | 13100 |

### Database

```sh
mise run migrate                    # Apply all pending migrations
mise run migrate-down               # Rollback one migration
mise run seed                       # Re-insert sample data
```

## Frontend Commands

All frontend commands run from `webapp/` or from the root using `pnpm -C webapp <command>`.

```sh
pnpm install                        # Install dependencies
pnpm dev                            # Dev server (port 13001)
pnpm test                           # Run tests (vitest)
pnpm lint                           # Lint (oxlint)
pnpm fmt                            # Format check (oxfmt)
pnpm fmt:fix                        # Apply formatting
pnpm check                          # Type check + lint + format check
pnpm check:fix                      # Type check + lint fix + format apply
pnpm build                          # Production build
pnpm preview                        # Preview production build
```

## Import Rules

- Vite config: `import { defineConfig } from 'vite';`
- Test utilities: `import { describe, expect, it, vi } from 'vitest';`

## Review Checklist for Agents

- [ ] Run `pnpm install` after pulling remote changes and before starting frontend work.
- [ ] Run `pnpm check` and `pnpm test` to validate frontend changes.
- [ ] Run `go test ./...` and `go tool golangci-lint run` to validate backend changes.
