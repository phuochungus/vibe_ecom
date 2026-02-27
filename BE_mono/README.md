# Golf Store Backend Monolith (BE_mono)

Monolithic Go backend for MVP phase.

## MySQL (real)

```bash
cd BE_mono
make mysql-up
```

Set env (or copy from `.env.example`):

```bash
export MYSQL_DSN='golf:golf@tcp(127.0.0.1:3307)/golf_store_mono?parseTime=true&charset=utf8mb4&loc=UTC'
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
- `MYSQL_DSN` is required at runtime.
- Data is persisted directly in MySQL tables (`users`, `products`, `orders`, `order_items`, `payment_transactions`, `notifications`, ...).
- Seed accounts:
  - `admin@golf.local / admin123`
  - `user@golf.local / user123`
