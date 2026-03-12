# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository overview
- This repo has 3 active apps; there is no root workspace. Run commands inside the relevant subdirectory.
- `BE_mono/` is the active backend monolith. `BE/` is legacy/reference material and is not the deployed backend.
- `FE_admin/` is the React/Vite admin SPA.
- `FE_customer/` is the React/Vite customer SPA.
- `render.yaml` defines the deployed topology: 1 Go API service, 2 static frontend sites, and 1 Postgres database.
- `tests/` stores shared Playwright artifacts and regression reports; the actual frontend specs live under each app's `tests/e2e/`.

## Common commands

### Backend (`BE_mono`)
- `cd BE_mono && make postgres-up`
- `cd BE_mono && make run`
- `cd BE_mono && make test`
- `cd BE_mono && go build -o bin/api ./cmd/api`
- `cd BE_mono && go test ./internal/modules/order/service -run TestName`
- `cd BE_mono && make postgres-down`
- Required local env comes from `BE_mono/.env.example`; `POSTGRES_DSN` is mandatory.
- Local API docs: `http://localhost:8080/docs`

### Admin frontend (`FE_admin`)
- `cd FE_admin && npm install`
- `cd FE_admin && npm run dev`
- `cd FE_admin && npm run build`
- `cd FE_admin && npm run lint`
- `cd FE_admin && npm run preview`
- `cd FE_admin && npm run test:e2e`
- `cd FE_admin && npm run test:e2e -- tests/e2e/admin-sprint1.spec.ts`
- `cd FE_admin && npm run test:e2e -- -g "admin can create a product"`
- `FE_admin` Playwright starts its own dev server on `127.0.0.1:4173`.

### Customer frontend (`FE_customer`)
- `cd FE_customer && npm install`
- `cd FE_customer && npm run dev`
- `cd FE_customer && npm run build`
- `cd FE_customer && npm run lint`
- `cd FE_customer && npm run preview`
- `cd FE_customer && npm run test:e2e`
- `cd FE_customer && npm run test:e2e -- tests/e2e/customer-sprint1.spec.ts`
- `cd FE_customer && npm run test:e2e -- -g "customer can complete cod checkout"`
- `FE_customer` dev server runs on `http://localhost:3001`.
- `FE_customer` Playwright does not start the dev server for you; run the app first, then run `npm run test:e2e`.

## Local development wiring
- Both frontends proxy `/api` to `http://localhost:8080` in local Vite config.
- In deployed environments, both frontends build the API base URL from `VITE_API_BASE_URL` or `VITE_API_HOST`.
- The backend serves Swagger UI at `/docs` and reads the OpenAPI file from `docs/FE/openapi.yaml` unless `OPENAPI_SPEC_PATH` is set.

## Backend architecture (`BE_mono`)
- Entrypoint: `cmd/api/main.go` loads env, builds config, bootstraps the app, and runs the HTTP server with graceful shutdown.
- App wiring: `internal/app/bootstrap/app.go` is the composition root. It opens Postgres, initializes schema, seeds demo data, constructs repositories/services/handlers, and passes them into the router.
- HTTP routing: `internal/app/server/router.go` defines:
  - health endpoints (`/healthz`, `/readyz`)
  - docs endpoints (`/docs`, `/openapi.yaml`)
  - `/api/v1` public routes
  - authenticated user routes protected by JWT middleware
  - admin routes protected by both auth and `RoleAdmin`
- Domain code lives under `internal/modules/<domain>/{http,service,repository,dto}`.
- Shared cross-cutting code lives under `internal/shared/` (middleware, errors, response helpers, utilities).
- Platform/infrastructure code lives under `internal/platform/` (config, db, entities, observability, http server).
- Schema creation and demo seeding happen at startup; local Postgres is defined in `docker-compose.postgres.yml` on port `5433`.
- When adding or changing backend modules, update both the bootstrap wiring in `internal/app/bootstrap/app.go` and route registration in `internal/app/server/router.go`.

## Frontend architecture
- Both frontends are independent React 19 + Vite + TypeScript SPAs with:
  - `createBrowserRouter`
  - TanStack Query for server state
  - Axios clients in `src/lib/api.ts`
  - localStorage-backed auth providers in `src/lib/auth.tsx`
  - `@/` aliased to `src`
- Both `src/lib/api.ts` implementations:
  - build the API base URL from env
  - attach `Authorization: Bearer ...`
  - attempt `/auth/refresh` on 401
  - clear stored tokens and redirect to `/login` if refresh fails
- Both frontends expect the backend response envelope shape `{ success, data }`.

### `FE_admin`
- Uses Tailwind/shadcn-style UI primitives and an authenticated admin shell.
- App providers are wired in `src/main.tsx`.
- `src/router.tsx` protects an entire authenticated route subtree and wraps it in `AdminLayout`.
- `src/lib/auth.tsx` allows only `ADMIN` users to log in.
- Service modules are thin wrappers around admin endpoints such as `/api/v1/admin/products`, `/api/v1/admin/orders`, and `/api/v1/admin/revenue`.

### `FE_customer`
- Uses Ant Design with centralized theme/locale setup in `src/main.tsx`.
- App-wide providers are wired in `src/App.tsx`.
- `src/router.tsx` keeps the storefront layout public, then wraps checkout/orders/notifications/profile routes individually with `ProtectedRoute`.
- `src/lib/auth.tsx` allows only `USER` users to log in and clears `shopping_cart` on logout.
- `src/lib/cart.tsx` adds a second client-side state layer: a localStorage-backed cart with stock/status checks.

## Testing notes
- Backend tests are standard Go `*_test.go` files; some service tests are integration-oriented and may require Postgres.
- Frontend Playwright suites live in `FE_admin/tests/e2e/` and `FE_customer/tests/e2e/`.
- Those Playwright specs stub API responses with `page.route(...)`, so they primarily validate frontend flows without a live backend.
- Root `tests/` contains Playwright evidence output and regression reports.

## Useful local credentials
- Admin seed account: `admin@golf.local / 123456`
- Customer seed account: `user@golf.local / 123456`
