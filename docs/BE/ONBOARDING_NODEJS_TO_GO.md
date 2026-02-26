# ONBOARDING: Node.js Dev -> Go Backend

## 1. Mục tiêu tài liệu
Bạn là dev Node.js và muốn hiểu codebase Go hiện tại nhanh nhất để bắt đầu làm task.

Mục tiêu thực tế:
1. Chạy được BE local.
2. Đọc được luồng request end-to-end.
3. Tự thêm được 1 endpoint mới ở API Gateway.
4. Hiểu cách service nội bộ giao tiếp qua broker (không HTTP service-to-service).

---

## 2. Map tư duy Node.js sang Go (rút gọn)

| Node.js | Go trong dự án này |
|---|---|
| Express app | Gin router |
| Route handler `(req,res)` | Method nhận `*gin.Context` |
| Class service | `struct` + method |
| Interface TS | `interface` Go |
| `throw`/`try-catch` | `if err != nil { ... }` |
| Async function | Goroutine `go fn()` |
| Env config (`process.env`) | `os.LookupEnv` |
| DTO/Type | Struct có tag JSON |

Ví dụ nhanh:
- Node: `res.status(202).json(...)`
- Go: `c.JSON(http.StatusAccepted, gin.H{...})`

---

## 3. Kiến trúc BE hiện tại (skeleton)

1. `api-gateway` là điểm vào HTTP duy nhất cho web client + webhook ngoài hệ thống.
2. Service nội bộ (`user/product/order/payment/notification`) chưa mở API nghiệp vụ public.
3. Giao tiếp nội bộ được chuẩn bị theo event/command qua broker (Kafka/RabbitMQ), không gọi HTTP trực tiếp giữa service.

Tham chiếu:
- [TECHSTACK.md](D:/golf_store/docs/BE/TECHSTACK.md)

---

## 4. Cây thư mục cần đọc trước

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
```

Ý nghĩa:
1. `cmd/*/main.go`: điểm khởi động mỗi service (như `server.ts` trong Node).
2. `internal/platform/*`: phần dùng chung (config, HTTP server wrapper, message envelope, correlation id).
3. `internal/services/*`: logic theo service/domain.

---

## 5. Cách đọc code nhanh nhất (đúng thứ tự)

### Bước 1: vào gateway trước
1. [cmd/api-gateway/main.go](D:/golf_store/BE/cmd/api-gateway/main.go)
2. [internal/services/gateway/app/app.go](D:/golf_store/BE/internal/services/gateway/app/app.go)
3. [internal/services/gateway/http/router.go](D:/golf_store/BE/internal/services/gateway/http/router.go)
4. [internal/services/gateway/http/handler.go](D:/golf_store/BE/internal/services/gateway/http/handler.go)

Bạn sẽ thấy flow:
`main -> app -> router -> handler -> publisher`.

### Bước 2: hiểu phần cross-cutting
1. [internal/platform/config/config.go](D:/golf_store/BE/internal/platform/config/config.go)
2. [internal/platform/observability/correlation.go](D:/golf_store/BE/internal/platform/observability/correlation.go)
3. [internal/platform/messaging/envelope.go](D:/golf_store/BE/internal/platform/messaging/envelope.go)

### Bước 3: nhìn 1 domain service mẫu
1. [cmd/user-service/main.go](D:/golf_store/BE/cmd/user-service/main.go)
2. [internal/services/user/app/app.go](D:/golf_store/BE/internal/services/user/app/app.go)
3. [internal/services/user/messaging/consumer.go](D:/golf_store/BE/internal/services/user/messaging/consumer.go)

---

## 6. Endpoint hiện có ở Gateway (để test)

Từ [router.go](D:/golf_store/BE/internal/services/gateway/http/router.go):

1. `GET /healthz`
2. `GET /readyz`
3. `POST /api/v1/auth/login`
4. `GET /api/v1/products`
5. `POST /api/v1/orders`
6. `GET /api/v1/orders/:orderCode/tracking`
7. `GET /api/v1/notifications`
8. `POST /api/v1/webhooks/payment/:provider`
9. `GET /api/v1/admin/orders`
10. `GET /api/v1/admin/revenue/summary`

Lưu ý:
- Nhiều endpoint đang là skeleton trả dữ liệu mẫu hoặc `202 Accepted`.

---

## 7. Chạy local như một dev Node.js mới qua Go

### Prerequisite
1. Go `>=1.22`
2. Docker + Docker Compose

### Chạy test nhanh
```powershell
cd D:\golf_store\BE
go test ./...
```

### Chạy riêng gateway
```powershell
cd D:\golf_store\BE
go run ./cmd/api-gateway
```

### Test endpoint
```powershell
curl http://localhost:8080/healthz
curl http://localhost:8080/api/v1/products
```

### Chạy full stack bằng Docker
```powershell
cd D:\golf_store\BE
docker compose up --build
```

---

## 8. Go syntax tối thiểu bạn cần biết cho dự án này

1. Struct:
```go
type User struct {
    ID string `json:"id"`
}
```

2. Method:
```go
func (h *Handler) Health(c *gin.Context) { ... }
```

3. Error-first:
```go
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
    return
}
```

4. Goroutine:
```go
go consumer.Start(ctx)
```

5. Context để shutdown graceful:
```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
```

---

## 9. Bài tập đầu tay (khuyên làm ngay)

### Bài tập A: thêm endpoint ping ở gateway
1. Thêm route `GET /api/v1/ping` trong [router.go](D:/golf_store/BE/internal/services/gateway/http/router.go).
2. Thêm handler trả:
```json
{"message":"pong"}
```
3. Chạy `go test ./...` rồi `go run ./cmd/api-gateway`.

### Bài tập B: thêm command mới
1. Trong handler, thêm endpoint `POST /api/v1/orders/:orderCode/cancel`.
2. Publish command type `order.cancel.requested`.
3. Trả `202 Accepted` kèm `commandId`, `correlationId`.

Mục tiêu của 2 bài tập:
- Nắm pattern coding của codebase hiện tại.
- Không cần hiểu sâu broker implementation ngay vòng đầu.

---

## 10. Quy tắc quan trọng khi code trong repo này

1. Client chỉ gọi HTTP vào `api-gateway`.
2. Không tạo HTTP call nghiệp vụ giữa services.
3. Giao tiếp nội bộ qua message broker (Kafka/RabbitMQ).
4. Mọi flow quan trọng giữ `correlation id` xuyên suốt.
5. Ưu tiên idempotent cho consumer/webhook.

---

## 11. Nếu bạn muốn onboard nhanh hơn nữa

Lộ trình 2 ngày:
1. Ngày 1 sáng: đọc mục `5`, chạy local, test route.
2. Ngày 1 chiều: làm Bài tập A.
3. Ngày 2 sáng: làm Bài tập B.
4. Ngày 2 chiều: bắt đầu implement thật một use case nhỏ (ví dụ login request publish qua RabbitMQ).

Sau lộ trình này, bạn đủ tự tin nhận task BE Go ở sprint hiện tại.
