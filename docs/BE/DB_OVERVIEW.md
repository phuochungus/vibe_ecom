# DB Overview - Golf Store (MySQL 8)

## 1. Mục tiêu
Tài liệu này định nghĩa chuẩn thiết kế CSDL backend cho hệ thống golf store, bám theo `docs/BRD.md` phiên bản `1.2`.

## 2. Phạm vi dữ liệu
- Quản lý người dùng và phân quyền (`USER`, `ADMIN`).
- Quản lý danh mục sản phẩm, tồn kho, trạng thái bán.
- Quản lý đơn hàng, thanh toán tự động, tracking giao vận.
- Quản lý thông báo theo sự kiện nghiệp vụ.
- Lưu audit log cho thực thể trọng yếu.

## 3. Chuẩn kỹ thuật bắt buộc
- DB engine: `MySQL 8`.
- Storage engine: `InnoDB`.
- Charset: `utf8mb4`.
- Collation khuyến nghị: `utf8mb4_0900_ai_ci`.
- Múi giờ lưu trữ: `UTC`.
- Tiền tệ mặc định: `VND` (`currency_code = 'VND'`).

## 4. Quy ước ID và thời gian
- Tất cả khóa chính: `id BINARY(16)` (UUID).
- Ở tầng ứng dụng sử dụng:
1. `UUID_TO_BIN(:uuid, 1)` khi ghi.
2. `BIN_TO_UUID(id, 1)` khi đọc.
- Cột thời gian chuẩn cho hầu hết bảng:
1. `created_at DATETIME(3) NOT NULL`
2. `updated_at DATETIME(3) NOT NULL`
- Soft delete áp dụng cho bảng danh mục: dùng `deleted_at DATETIME(3) NULL`.

## 5. Quy ước đặt tên
- Tên bảng/cột: `snake_case`.
- Tên FK: `fk_<table>__<ref_table>`.
- Tên unique index: `uk_<table>__<columns>`.
- Tên non-unique index: `idx_<table>__<columns>`.

## 6. Quy ước dữ liệu
- Số tiền: `DECIMAL(18,2)`.
- Số lượng: `INT UNSIGNED`.
- Trạng thái nghiệp vụ: dùng `ENUM` (chuẩn hóa trong data dictionary).
- Metadata linh hoạt: dùng `JSON` (chỉ khi cần).
- Dữ liệu PII (email, phone): bắt buộc unique và kiểm soát truy cập ở tầng ứng dụng.

## 7. Assumptions
- Single-tenant, không multi-warehouse trong Sprint-01.
- Không có loyalty/point ở phiên bản hiện tại.
- Không triển khai migration SQL trong tài liệu này (doc-only).

## 8. Danh sách bảng chính
1. `users`
2. `user_addresses`
3. `products`
4. `orders`
5. `order_items`
6. `order_status_history`
7. `payment_transactions`
8. `shipment_tracking_events`
9. `notifications`
10. `audit_logs`
