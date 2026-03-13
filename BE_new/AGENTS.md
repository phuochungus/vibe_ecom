# Repository Guidelines

## Project Structure & Module Organization
This repository is a Go backend service. Runtime entrypoints live in `cmd/`: `cmd/api` starts the HTTP server and `cmd/migrate` applies database migrations. Core application code is under `internal/`, split by responsibility (`config`, `database`, `handlers`, `models`, `routes`, `services`, `repositories`, `middleware`). SQL migrations live in `migrations/` with paired files such as `000001_init.up.sql` and `000001_init.down.sql`. Keep generated API docs in `docs/`.

## Build, Test, and Development Commands
Use the `Makefile` for common local tasks:

- `make infra-up`: start PostgreSQL and MinIO from `docker-compose.yaml`.
- `make run`: run the API with `go run ./cmd/api`.
- `make migrate-up`: apply pending migrations.
- `make migrate-down`: roll back the latest migration batch.
- `go test ./...`: run the full Go test suite.

Use `make infra-down` to stop containers and remove volumes when resetting local state.

## Coding Style & Naming Conventions
Follow standard Go conventions. Format all edits with `gofmt` and keep imports organized. Use tabs for indentation, short lowercase package names, and descriptive exported identifiers in `PascalCase`. Match existing file naming patterns such as `user_role.go` and `order_item.go`. Keep HTTP wiring in `routes`, request handling in `handlers`, and persistence logic out of route setup.

## Testing Guidelines
Add tests alongside the code they cover using `*_test.go` files and the standard `testing` package. Prefer table-driven tests for config loading, handlers, and database-facing logic. There are currently no committed tests, so new features and bug fixes should include coverage where behavior can regress. Run `go test ./...` before opening a PR.

## Commit & Pull Request Guidelines
Recent history uses Conventional Commit prefixes such as `feat:`, `fix:`, and `refactor:`. Keep commit subjects imperative and scoped to one change. PRs should explain the behavior change, note any schema or `.env` updates, and include verification steps (for example, `make migrate-up` and `go test ./...`). Attach request/response samples when API behavior changes.

## Configuration & Infrastructure Tips
Local configuration comes from `.env`, loaded automatically at startup. Database settings must include `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSL_MODE`, and `DB_TIMEZONE`. The server currently reads `APP_PORT`, so keep that key populated when editing environment files.
