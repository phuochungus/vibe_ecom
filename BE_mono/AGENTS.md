# Repository Guidelines

## Project Structure & Module Organization
`cmd/api` contains the API entrypoint. Application wiring lives in `internal/app`, while reusable infrastructure sits in `internal/platform` (`config`, `db`, `entities`, `observability`). Business code is organized by domain under `internal/modules/<domain>/{http,service,repository,dto}`. Cross-cutting helpers live in `internal/shared` (`middleware`, `response`, `errors`, `utils`). Database schema changes go in `migrations/`; local MySQL orchestration is defined in `docker-compose.mysql.yml`.

## Build, Test, and Development Commands
Use the Makefile first:

- `make mysql-up`: start the local MySQL container on the configured port.
- `make run`: run the API with `go run ./cmd/api`.
- `make test`: run the full Go test suite with `go test ./...`.
- `make mysql-down`: stop MySQL and remove the local volume.

Before `make run`, set `MYSQL_DSN` and auth settings from `.env.example`. `MYSQL_DSN` is required at runtime.

## Coding Style & Naming Conventions
Keep code `gofmt`-clean; no separate linter config is checked in. Follow standard Go formatting, grouped imports, and lowercase package names. Use `snake_case.go` filenames such as `order_status_history.go`. Keep HTTP handlers thin, place business rules in `service`, and isolate GORM/data access in `repository`. Exported identifiers use `PascalCase`; local variables use `camelCase`.

## Testing Guidelines
Place tests beside the package they cover in `*_test.go` files and use `TestXxx` names. Run `make test` before opening a PR. Existing service tests in `internal/modules/*/service` are integration-oriented and may require a configured MySQL database; guard those cases explicitly with skips when the database is unavailable. Add tests for new service logic, request validation, and bootstrap/config edge cases where practical.

## Commit & Pull Request Guidelines
Recent history favors Conventional Commit-style subjects like `feat: ...` and `refactor: ...`. Keep the type lowercase, write the subject in the imperative, and make it specific to the changed module, for example `feat: add payment status update endpoint`.

PRs should include a short summary, affected modules, migration or env var changes, and the commands you ran to verify the change. For API or schema updates, include a sample request/response, route list, or migration note so reviewers can validate behavior quickly.

## Configuration & Data
Never commit real secrets in `.env`. Keep machine-specific values such as `OPENAPI_SPEC_PATH` local only. Add schema changes as new numbered files in `migrations/` instead of rewriting applied migrations.
