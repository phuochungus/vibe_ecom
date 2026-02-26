# API Contract - FE <-> API Gateway

## 1. Scope
- Version: `v1`
- Public base path: `/api/v1`
- Chỉ `api-gateway` public ra ngoài.
- FE chỉ gọi API Gateway, không gọi trực tiếp service nội bộ.

## 2. Chuẩn chung
### 2.1 Headers
- `Authorization: Bearer <access_token>` (trừ endpoint login/refresh).
- `Content-Type: application/json`.
- `X-Request-Id` (optional, khuyến nghị FE gửi).
- `Idempotency-Key` (bắt buộc cho create order/payment).

### 2.2 Dữ liệu chuẩn
- ID: UUID string.
- Timestamp: ISO-8601 UTC (`2026-02-26T15:30:00.000Z`).
- Amount: decimal string (`"1590000.00"`).

### 2.3 Envelope response
```json
{
  "success": true,
  "data": {},
  "meta": {},
  "request_id": "5f1e...",
  "timestamp": "2026-02-26T15:30:00.000Z"
}
```

### 2.4 Error response
```json
{
  "success": false,
  "error": {
    "code": "ORDER_NOT_FOUND",
    "message": "Order not found",
    "details": []
  },
  "request_id": "5f1e...",
  "timestamp": "2026-02-26T15:30:00.000Z"
}
```

## 3. RBAC Matrix
- `USER`: auth, products, orders của chính mình, payments của chính mình, notifications.
- `ADMIN`: toàn bộ `USER` + quản trị products, orders, revenue.

## 4. Authentication APIs
### 4.1 POST `/auth/login`
- Auth: Public
- Request:
```json
{ "identifier": "user@email.com", "password": "***" }
```
- Response `200`:
```json
{
  "access_token": "...",
  "refresh_token": "...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": {
    "id": "uuid",
    "role": "USER",
    "full_name": "Nguyen Van A",
    "email": "user@email.com",
    "phone": "090...",
    "status": "ACTIVE"
  }
}
```
- Errors: `401 INVALID_CREDENTIALS`, `423 ACCOUNT_LOCKED`.

### 4.2 POST `/auth/refresh`
- Auth: Public
- Request: `{ "refresh_token": "..." }`
- Response `200`: token mới.

### 4.3 POST `/auth/logout`
- Auth: Bearer
- Response `204`.

### 4.4 GET `/me`
- Auth: Bearer
- Response `200`: profile hiện tại.

## 5. Product APIs
### 5.1 GET `/products`
- Auth: Optional
- Query:
1. `q`
2. `status=ACTIVE|INACTIVE|DISCONTINUED` (admin)
3. `min_price`, `max_price`
4. `page`, `page_size`
5. `sort=created_at|price|name`
6. `order=asc|desc`
- Response `200`:
```json
{
  "items": [
    {
      "id": "uuid",
      "sku": "CLUB-001",
      "name": "Driver X",
      "price": "1590000.00",
      "stock": 10,
      "status": "ACTIVE",
      "image_url": "https://..."
    }
  ],
  "pagination": { "page": 1, "page_size": 20, "total": 120, "total_pages": 6 }
}
```

### 5.2 GET `/products/{product_id}`
- Auth: Optional
- Response `200`: chi tiết sản phẩm.

## 6. Order APIs (User)
### 6.1 POST `/orders`
- Auth: USER
- Headers: `Idempotency-Key` bắt buộc
- Request:
```json
{
  "items": [
    { "product_id": "uuid", "quantity": 1 }
  ],
  "shipping_address": {
    "recipient_name": "Nguyen Van A",
    "recipient_phone": "090...",
    "line1": "...",
    "line2": "...",
    "ward": "...",
    "district": "...",
    "city": "HCM",
    "province": "HCM",
    "postal_code": "700000",
    "country_code": "VN"
  },
  "customer_note": "..."
}
```
- Response `201`:
```json
{
  "id": "uuid",
  "order_code": "ORD-20260226-0001",
  "order_status": "PENDING_PAYMENT",
  "payment_status": "UNPAID",
  "subtotal_amount": "1590000.00",
  "discount_amount": "0.00",
  "shipping_fee": "30000.00",
  "total_amount": "1620000.00",
  "payment_due_at": "2026-02-26T16:00:00.000Z",
  "placed_at": "2026-02-26T15:30:00.000Z"
}
```
- Errors: `400 INVALID_ITEMS`, `409 OUT_OF_STOCK`.

### 6.2 GET `/orders`
- Auth: USER
- Query: `status`, `from`, `to`, `page`, `page_size`
- Response `200`: danh sách đơn của user.

### 6.3 GET `/orders/{order_id}`
- Auth: USER
- Response `200`: chi tiết đơn + items + payment summary.

### 6.4 POST `/orders/{order_id}/cancel`
- Auth: USER
- Request: `{ "reason": "..." }`
- Response `200`: đơn đã chuyển `CANCELLED` (nếu hợp lệ rule).
- Errors: `409 ORDER_CANNOT_CANCEL`.

### 6.5 GET `/orders/{order_id}/tracking`
- Auth: USER
- Response `200`:
```json
{
  "order_id": "uuid",
  "current_status": "SHIPPING",
  "timeline": [
    {
      "status": "PENDING_PAYMENT",
      "source_type": "SYSTEM",
      "occurred_at": "2026-02-26T15:30:00.000Z",
      "description": "Order created"
    }
  ]
}
```

## 7. Payment APIs (User)
### 7.1 POST `/orders/{order_id}/payments`
- Auth: USER
- Headers: `Idempotency-Key` bắt buộc
- Request:
```json
{
  "provider": "PAYOS",
  "return_url": "https://fe.app/payment/return",
  "cancel_url": "https://fe.app/payment/cancel"
}
```
- Response `201`:
```json
{
  "payment_id": "uuid",
  "status": "PENDING",
  "checkout_url": "https://payment-provider/..."
}
```

### 7.2 GET `/orders/{order_id}/payments`
- Auth: USER
- Response `200`: danh sách payment transaction của đơn.

## 8. Notification APIs (User)
### 8.1 GET `/notifications`
- Auth: USER
- Query: `status=PENDING|SENT|FAILED`, `page`, `page_size`
- Response `200`: danh sách thông báo.

### 8.2 PATCH `/notifications/{notification_id}/read`
- Auth: USER
- Response `200`: cập nhật trạng thái đã đọc.

### 8.3 PATCH `/notifications/read-all`
- Auth: USER
- Response `200`: đánh dấu tất cả đã đọc.

## 9. Admin APIs
### 9.1 Products (Admin)
- `POST /admin/products`
- `PATCH /admin/products/{product_id}`
- `DELETE /admin/products/{product_id}` (soft delete)

Create request mẫu:
```json
{
  "sku": "CLUB-001",
  "name": "Driver X",
  "description": "...",
  "price": "1590000.00",
  "stock": 100,
  "status": "ACTIVE",
  "image_url": "https://..."
}
```

### 9.2 Orders (Admin)
- `GET /admin/orders`
- `GET /admin/orders/{order_id}`
- `PATCH /admin/orders/{order_id}/status`

Status update request:
```json
{
  "to_status": "PROCESSING",
  "reason": "verified payment"
}
```
- Rule: chỉ chuyển trạng thái theo state machine chuẩn.

### 9.3 Revenue (Admin)
- `GET /admin/revenue/summary?from=...&to=...&group_by=day|month|year`
- `GET /admin/revenue/orders?from=...&to=...&page=1&page_size=20`

Summary response mẫu:
```json
{
  "from": "2026-02-01T00:00:00.000Z",
  "to": "2026-02-29T00:00:00.000Z",
  "gross_revenue": "50000000.00",
  "refund_amount": "2000000.00",
  "net_revenue": "48000000.00",
  "completed_orders": 320,
  "payment_success_rate": "0.93"
}
```

## 10. Webhook Endpoint (External -> Gateway)
### 10.1 POST `/webhooks/payments/{provider}`
- Caller: payment provider (không phải FE).
- Qua API Gateway, gateway validate signature/basic shape rồi publish vào broker.
- Response: `202 Accepted`.

## 11. Status/Enum chuẩn
- `user_role`: `USER`, `ADMIN`
- `product_status`: `ACTIVE`, `INACTIVE`, `DISCONTINUED`
- `order_status`: `NEW`, `PENDING_PAYMENT`, `PAID`, `PROCESSING`, `SHIPPING`, `COMPLETED`, `CANCELLED`
- `payment_status`: `UNPAID`, `PAID`, `FAILED`, `REFUNDED`
- `payment_txn_state`: `PENDING`, `SUCCESS`, `FAILED`, `CANCELLED`
- `notification_status`: `PENDING`, `SENT`, `FAILED`

## 12. Error code đề xuất
- `UNAUTHORIZED`, `FORBIDDEN`
- `INVALID_REQUEST`, `VALIDATION_ERROR`
- `RESOURCE_NOT_FOUND`
- `OUT_OF_STOCK`
- `ORDER_CANNOT_CANCEL`
- `PAYMENT_TIMEOUT`, `PAYMENT_FAILED`, `PAYMENT_DUPLICATE`
- `RATE_LIMITED`
- `INTERNAL_ERROR`

## 13. OpenAPI next step
- Từ contract này, tạo file `openapi.yaml` làm single source of truth.
- FE dùng codegen client từ OpenAPI để giảm lỗi integrate.
