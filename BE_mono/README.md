# Golf Store Backend Monolith (BE_mono)

Monolithic Go backend for MVP phase.

## PostgreSQL (real)

```bash
cd BE_mono
make postgres-up
```

Set env (or copy from `.env.example`):

```bash
export POSTGRES_DSN='host=127.0.0.1 user=golf password=golf dbname=golf_store_mono port=5433 sslmode=disable TimeZone=UTC'
export JWT_SECRET='replace_with_strong_secret'
export JWT_ISSUER='be-mono'
export JWT_ACCESS_TTL_MINUTES='15'
export JWT_REFRESH_TTL_MINUTES='10080'
# optional: override OpenAPI file path (default: docs/FE/openapi.yaml or ../docs/FE/openapi.yaml)
export OPENAPI_SPEC_PATH='/Users/imacvip/golf_store/docs/FE/openapi.yaml'
```

## Run

```bash
cd BE_mono
go run ./cmd/api
```

Open docs:

```bash
open http://localhost:8080/docs
```

## Test

```bash
cd BE_mono
go test ./...
```

## Notes

- New source code lives in `BE_mono/` only.
- Existing microservice skeleton in `BE/` is untouched.
- `POSTGRES_DSN` is required at runtime.
- Data is persisted directly in PostgreSQL tables (`users`, `products`, `orders`, `order_items`, `payment_transactions`, `notifications`, ...).
- Seed accounts:
  - `admin@golf.local / 123456`
  - `user@golf.local / 123456`
