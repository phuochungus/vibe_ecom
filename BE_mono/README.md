# Golf Store Backend Monolith (BE_mono)

Monolithic Go backend for MVP phase.

## Local Infra

```bash
cd BE_mono
make infra-up
```

Set env (or copy from `.env.example`):

```bash
export POSTGRES_DSN='host=127.0.0.1 user=golf password=golf dbname=golf_store_mono port=5433 sslmode=disable TimeZone=UTC'
export MINIO_ENABLED='true'
export MINIO_ENDPOINT='127.0.0.1:9000'
export MINIO_PUBLIC_BASE_URL='http://127.0.0.1:9000'
export MINIO_ACCESS_KEY='minioadmin'
export MINIO_SECRET_KEY='minioadmin'
export MINIO_BUCKET='golf-store'
export MINIO_USE_SSL='false'
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
make run
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
- Local Docker services are defined in `docker-compose.yaml`.
- Data is persisted directly in PostgreSQL tables (`users`, `products`, `orders`, `order_items`, `payment_transactions`, `notifications`, ...).
- Admin product image uploads are available at `POST /api/v1/admin/products/upload-image` when MinIO is enabled. Send the image as multipart form-data in the `file` field.
- Seed accounts:
  - `admin@golf.local / 123456`
  - `user@golf.local / 123456`
