# Monolithic Architecture - MVP (BE_mono)

## 1. Thông tin tài liệu
- Mã tài liệu: `ARCH-GOLF-MONO-001`
- Phiên bản: `1.0`
- Ngày ban hành: `27/02/2026`
- Trạng thái: `Draft for implementation`
- Phạm vi áp dụng: backend mới cho MVP tại thư mục `BE_mono/`

## 2. Mục tiêu
- Triển khai backend Go theo kiến trúc monolithic để rút ngắn thời gian ra MVP.
- Tái sử dụng đầy đủ nghiệp vụ đã chốt trong `docs/BRD.md`, `docs/business.md`, `docs/FE/openapi.yaml`.
- Không thay đổi và không can thiệp code hiện tại trong `BE/`.

## 3. Quyết định kiến trúc
- Kiến trúc chính: `Modular Monolith`.
- Một process HTTP duy nhất cho toàn bộ domain nghiệp vụ.
- Một database logic chung (PostgreSQL 16) cho MVP.
- Giao tiếp nội bộ giữa domain bằng function call + transaction, không tách service mạng nội bộ.
- Bất đồng bộ (nếu cần) chạy qua outbox/job trong cùng codebase monolith.

## 4. Vì sao chọn monolith cho giai đoạn MVP
- Tối ưu tốc độ phát triển: giảm chi phí vận hành nhiều service.
- Debug đơn giản: trace request trong một process.
- Dễ đảm bảo transaction xuyên domain ở các luồng tạo đơn/thanh toán.
- Vẫn giữ ranh giới module rõ ràng để có thể tách microservice sau này.

## 5. Nguyên tắc thiết kế bắt buộc
- Source code mới đặt tại `BE_mono/`.
- Không sửa, không di chuyển, không xóa file trong `BE/`.
- API công khai vẫn theo chuẩn `v1` và bám contract tại `docs/FE/openapi.yaml`.
- Domain boundary rõ theo module: `auth`, `product`, `order`, `payment`, `notification`, `reporting`.
- Mọi thay đổi trạng thái đơn/thanh toán phải có audit log.
- Rule nghiệp vụ bám `BR-*` trong BRD, không tự diễn giải khác.

## 6. Cấu trúc thư mục đề xuất cho `BE_mono`
```text
BE_mono/
  cmd/
    api/
      main.go
  internal/
    app/
      bootstrap/
      server/
    platform/
      config/
      db/
      cache/
      queue/              # optional cho async nội bộ
      observability/
    shared/
      dto/
      errors/
      middleware/
      response/
      utils/
    modules/
      auth/
        http/
        service/
        repository/
        model/
      product/
        http/
        service/
        repository/
        model/
      order/
        http/
        service/
        repository/
        model/
      payment/
        http/
        service/
        repository/
        model/
      notification/
        http/
        service/
        repository/
        model/
      reporting/
        http/
        service/
        repository/
  migrations/
  scripts/
  docs/
  go.mod
```

## 7. Luồng xử lý chuẩn trong monolith
1. HTTP request vào router (`/api/v1/...`).
2. Middleware xử lý auth, correlation id, request id, validation cơ bản.
3. Handler map request DTO -> service input.
4. Service xử lý business rule + transaction.
5. Repository thao tác DB.
6. Service ghi `audit_logs` và phát sự kiện nội bộ (nếu có).
7. Handler trả response envelope thống nhất cho FE.

## 8. Mapping module theo nghiệp vụ MVP
- `auth`: login/refresh/logout/me, lock account sau 5 lần sai.
- `product`: danh sách sản phẩm, chi tiết sản phẩm, admin CRUD sản phẩm.
- `order`: tạo đơn, danh sách đơn, chi tiết đơn, hủy đơn, tracking.
- `payment`: tạo payment, nhận webhook callback, idempotency xử lý callback.
- `notification`: danh sách thông báo, read/read-all, dedupe theo event_key.
- `reporting`: revenue summary, revenue orders breakdown.

## 9. Dữ liệu và transaction
- Schema bám `docs/BE/DATA_DICTIONARY.md`.
- Các transaction bắt buộc:
1. Tạo order + order_items + status_history.
2. Thanh toán callback: payment_transactions + update orders + status_history + audit.
3. Hủy đơn hợp lệ + status_history + audit.
- Số tiền dùng `DECIMAL(18,2)`, timezone lưu UTC.

## 10. Async trong monolith (khuyến nghị)
- Dùng bảng outbox để lưu event sau commit transaction.
- Worker nội bộ poll outbox và xử lý:
1. gửi notification
2. đồng bộ read-model/report cache
- Mọi consumer nội bộ phải idempotent theo event id / idempotency key.

## 11. Quy chuẩn API cho FE
- Base path: `/api/v1`.
- Giữ nguyên header policy:
1. `Authorization: Bearer <token>`
2. `Idempotency-Key` cho create order/payment
3. `X-Request-Id` (optional)
- Giữ nguyên envelope response/error như `docs/FE/API_CONTRACT.md`.

## 12. Bảo mật và quan sát hệ thống
- Password hash bằng thuật toán phù hợp (vd: bcrypt/argon2).
- Không log dữ liệu nhạy cảm thanh toán.
- Structured logging (JSON) + correlation id toàn request.
- Metrics tối thiểu:
1. HTTP latency/error rate
2. payment callback success rate
3. job retry/failure

## 13. Kế hoạch triển khai BE_mono theo Sprint-01
### Phase A - Bootstrap
- Khởi tạo skeleton `BE_mono`.
- Thiết lập config, router, middleware, DB connection, migration.
- Dựng health/readiness endpoint.

### Phase B - Core nghiệp vụ
- Implement `auth`, `product`, `order`, `payment` theo AC hiện tại.
- Bổ sung `notification` và `reporting`.
- Hoàn thiện audit log + tracking timeline.

### Phase C - Hardening
- Unit test domain core.
- Integration test cho flow create order/payment callback/reporting.
- Kiểm thử contract bám `openapi.yaml`.

## 14. Definition of Done cho monolith MVP
- Pass test cho các flow `US-01` đến `US-07`.
- API tương thích `docs/FE/openapi.yaml`.
- Đảm bảo business rules trọng yếu:
1. `BR-ORD-03`, `BR-ORD-04`
2. `BR-PAY-01`..`BR-PAY-04`
3. `BR-REV-01`..`BR-REV-03`
- Có changelog và tài liệu cập nhật cho BE_mono.

## 15. Chuẩn bị cho tách microservices sau MVP
- Giữ boundary module rõ, không import chéo vòng tròn.
- Service layer không phụ thuộc trực tiếp transport.
- Event nội bộ dùng envelope thống nhất để dễ đổi sang broker thật.
- Ưu tiên interface cho repository và external adapter.

## 16. Kết luận
Kiến trúc `Modular Monolith` tại `BE_mono/` phù hợp mục tiêu ra MVP nhanh, giảm độ phức tạp vận hành, và vẫn giữ khả năng tách thành microservices ở giai đoạn mở rộng mà không làm ảnh hưởng code hiện tại trong `BE/`.
