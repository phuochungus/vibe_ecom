# Tech Stack - Backend Golf Store

## 1. Mục tiêu
Tài liệu này định nghĩa tech stack backend cho dự án golf store theo kiến trúc microservices, với nguyên tắc ingress duy nhất qua API Gateway.

## 2. Quyết định kiến trúc bắt buộc
1. Chỉ có `api-gateway` được public ra ngoài.
2. `Web client` chỉ gọi HTTP API vào `api-gateway`.
3. Webhook ngoài hệ thống (ví dụ cổng thanh toán) cũng chỉ được nhận tại `api-gateway`.
4. Các service backend **không giao tiếp với nhau qua HTTP/gRPC**.
5. Giao tiếp nội bộ service-to-service chỉ qua message broker (`Kafka`/`RabbitMQ`).

## 3. Core Stack
- Ngôn ngữ: `Golang` (khuyến nghị `>= 1.22`).
- HTTP framework: `Gin` (dùng cho `api-gateway` và endpoint nội bộ như health/metrics).
- Kiến trúc: `Microservices` + `Event-driven`.
- Cache/read model: `Redis`.
- Message broker:
1. `Kafka` (event streaming).
2. `RabbitMQ` (command queue, routing, retry).
- Containerization: `Docker`.
- Local orchestration: `Docker Compose`.

## 4. Thành phần hệ thống
### 4.1 Edge component
- `api-gateway`
- Vai trò:
1. Entry point HTTP duy nhất cho web client và webhook ngoài hệ thống.
2. AuthN/AuthZ ở tầng biên.
3. Validate request và phát command/event vào broker.
4. Trả response cho client theo mode đồng bộ (reply queue) hoặc bất đồng bộ (`202 Accepted`).

### 4.2 Domain microservices
1. `user-service`
- Quản lý user, auth domain, role (`USER`, `ADMIN`).
- Phát event user domain.

2. `product-service`
- Quản lý danh mục sản phẩm, giá, tồn kho, trạng thái bán.
- Phát event thay đổi sản phẩm/tồn kho.

3. `order-service`
- Quản lý đơn hàng và state machine trạng thái đơn.
- Phát event vòng đời đơn hàng.

4. `payment-service`
- Xử lý nghiệp vụ thanh toán, callback validation, idempotency.
- Phát event `payment.succeeded/failed/refunded`.

5. `notification-service`
- Nhận event nghiệp vụ và gửi thông báo.
- Retry + dedupe + trạng thái gửi.

## 5. Mô hình giao tiếp
### 5.1 Luồng vào hệ thống
- `Web client -> API Gateway (HTTP) -> RabbitMQ/Kafka -> Domain services`.
- `External webhook -> API Gateway (HTTP) -> RabbitMQ/Kafka -> Target service`.

### 5.2 Luồng giữa các service
- Chỉ dùng broker, không call HTTP trực tiếp.
- Service xử lý message theo consumer group/queue consumer với cơ chế idempotent.

### 5.3 Vai trò từng broker
- `RabbitMQ`:
1. Command/event ngắn hạn cần routing linh hoạt.
2. Request/reply qua `reply_to`/correlation_id cho luồng cần phản hồi nhanh.
3. Retry theo queue + DLX/DLQ.

- `Kafka`:
1. Event streaming, event log, replay.
2. Integration event giữa domain và analytics/reporting pipeline.
3. Consumer group scale-out cho xử lý bất đồng bộ.

### 5.4 Quy tắc naming
- RabbitMQ exchange: `cmd.<domain>`, `evt.<domain>`.
- RabbitMQ queue: `<service>.<purpose>.q`, DLQ: `<service>.<purpose>.dlq`.
- Kafka topic ví dụ:
1. `order.created`
2. `order.status.changed`
3. `payment.succeeded`
4. `payment.failed`
5. `payment.refunded`
6. `notification.requested`

### 5.5 Đảm bảo nhất quán
- Dùng `Outbox pattern` để publish event sau commit DB.
- Consumer bắt buộc idempotent (dedupe key/event id).
- Có retry + DLQ (RabbitMQ) và retry + DLQ topic/pattern (Kafka consumer side).

## 6. Redis Strategy
- Cache read-heavy: sản phẩm, danh sách sản phẩm, summary order/tracking.
- Cache token/session ngắn hạn và rate-limit key tại gateway.
- Cache idempotency key cho payment/webhook.
- Pattern:
1. `cache-aside` cho query.
2. TTL rõ ràng theo từng loại key.

## 7. Data Ownership
- Mỗi service sở hữu schema/domain data của chính nó.
- Không cho service truy cập trực tiếp bảng của service khác.
- Trao đổi dữ liệu xuyên service thông qua event/command qua broker hoặc read model đã materialize.

## 8. Testing Strategy
### 8.1 Unit test
- Tool: `go test`, `testify`, `gomock/mockery`.
- Scope: business logic, validation, state transition, idempotency handler.
- Yêu cầu: coverage tối thiểu `80%` cho domain core của từng service.

### 8.2 Integration test (service-level)
- Tool: `testcontainers-go`.
- Dependency test: `MySQL`, `Redis`, `Kafka`, `RabbitMQ`.
- Verify:
1. DB transaction + outbox.
2. Publish/consume Kafka đúng schema.
3. Publish/consume RabbitMQ đúng exchange/queue routing key.
4. Cache invalidation đúng rule.

### 8.3 Contract test
- Kafka event contract: schema versioned (`JSON Schema`/`Avro`).
- RabbitMQ message contract: schema versioned + routing key contract.
- Producer/consumer contract phải đảm bảo backward compatibility.

### 8.4 End-to-end test
- Chạy full stack bằng `docker compose`.
- E2E chỉ đi qua `api-gateway` (không gọi trực tiếp domain service).
- Kịch bản bắt buộc:
1. Login user.
2. Xem sản phẩm.
3. Tạo đơn.
4. Thanh toán + webhook qua gateway.
5. Tracking trạng thái đơn.
6. Nhận thông báo.
7. Kiểm tra báo cáo doanh thu.

### 8.5 Performance test
- Tool: `k6` (hoặc tương đương).
- Test tại gateway và critical path của Kafka/RabbitMQ.
- KPI tối thiểu:
1. API p95 dưới ngưỡng SLA nghiệp vụ.
2. Kafka consumer lag và RabbitMQ queue depth không vượt ngưỡng vận hành.

### 8.6 CI Quality Gate
- Bắt buộc pass: unit + integration + contract test.
- E2E chạy trên pipeline release hoặc nightly.
- Block merge nếu vi phạm schema compatibility hoặc test thất bại.

## 9. Docker & Docker Compose
### 9.1 Container list
- `api-gateway`
- `user-service`
- `product-service`
- `order-service`
- `payment-service`
- `notification-service`
- `redis`
- `kafka`
- `rabbitmq`
- `mysql` (hoặc cụm DB theo chiến lược triển khai)

### 9.2 Nguyên tắc network
- Chỉ expose port của `api-gateway` ra host/public.
- Các domain service chỉ mở trong network nội bộ compose/k8s.
- Webhook provider được route vào gateway endpoint.

### 9.3 Build/runtime
- Mỗi service dùng Dockerfile multi-stage.
- Healthcheck bắt buộc cho gateway và domain services.
- Graceful shutdown khi restart/deploy.

## 10. Config & Secrets
- Dùng biến môi trường (`.env`) theo service.
- Cấu hình tối thiểu:
1. DB connection
2. Redis URL
3. Kafka brokers + topic names
4. RabbitMQ URL + exchange/queue/routing key
5. JWT secret / signing keys
6. Payment provider keys
- Không commit secret thật vào repo.

## 11. Observability
- Logging dạng structured JSON.
- Correlation ID xuyên suốt HTTP request -> broker message chain.
- Metrics bắt buộc:
1. Gateway latency/error rate.
2. Kafka consumer lag + retry.
3. RabbitMQ queue depth + nack + DLQ count.
4. Redis hit/miss.
5. Payment callback success rate.

## 12. NFR kỹ thuật bắt buộc
- Idempotency cho webhook/payment callback và broker consumer.
- Không có HTTP service-to-service trong nghiệp vụ.
- Timeout + retry + circuit breaker cho thao tác I/O ngoài hệ thống.
- Backward-compatible API/event trong cùng major version.

## 13. Versioning
- API gateway versioning: `/api/v1/...`.
- Kafka topic schema versioning (`v1`, `v2`...).
- RabbitMQ message schema versioning (`v1`, `v2`...).
- Release version theo SemVer (`MAJOR.MINOR.PATCH`).

## 14. Deliverable liên quan
- Thiết kế CSDL chi tiết: `docs/BE/DATA_DICTIONARY.md`.
- Mapping business rule: `docs/BE/BUSINESS_RULE_MAPPING.md`.
- Checklist implement: `docs/BE/IMPLEMENTATION_CHECKLIST.md`.
