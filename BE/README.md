# Golf Store Backend Skeleton

Backend monorepo skeleton for Golf Store using:

- Go + Gin
- API Gateway as the only public HTTP entry
- Event-driven internals (Kafka + RabbitMQ placeholders)
- Redis cache
- MySQL 8
- Docker + Docker Compose

## Structure

```text
BE/
  cmd/
    api-gateway/
    user-service/
    product-service/
    order-service/
    payment-service/
    notification-service/
  internal/
    platform/
      config/
      httpserver/
      messaging/
      observability/
    services/
      gateway/
      user/
      product/
      order/
      payment/
      notification/
  docker/
    Dockerfile.service
  docker-compose.yml
  go.mod
```

## Quick start

1. Start dependencies and all services:

```bash
docker compose up --build
```

2. API Gateway is exposed at:

```text
http://localhost:8080
```

3. Health checks:

```text
GET /healthz
GET /readyz
```

4. Example gateway routes:

```text
POST /api/v1/auth/login
GET  /api/v1/products
POST /api/v1/orders
GET  /api/v1/orders/:orderCode/tracking
GET  /api/v1/notifications
POST /api/v1/webhooks/payment/:provider
```

## Notes

- Internal services only expose health/readiness in this skeleton.
- Service-to-service business communication is prepared via messaging stubs, not HTTP.
- This is a starter scaffold for iterative implementation.
