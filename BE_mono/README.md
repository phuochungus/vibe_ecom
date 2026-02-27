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
```

## Run

```bash
cd BE_mono
go run ./cmd/api
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
