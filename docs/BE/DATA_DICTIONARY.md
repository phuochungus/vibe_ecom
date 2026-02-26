# Data Dictionary - Golf Store (MySQL 8)

## 1. Quy ước chung
- PK chuẩn: `id BINARY(16)` (UUID).
- Timestamps: `DATETIME(3)` theo UTC.
- Số tiền: `DECIMAL(18,2)`.
- Soft delete dùng `deleted_at` cho bảng danh mục.
- Tất cả FK dùng `ON UPDATE RESTRICT`; `ON DELETE RESTRICT` cho dữ liệu giao dịch.

## 2. Enum domain chuẩn hóa
- `user_role`: `USER`, `ADMIN`
- `user_status`: `ACTIVE`, `LOCKED`, `DISABLED`
- `product_status`: `ACTIVE`, `INACTIVE`, `DISCONTINUED`
- `order_status`: `NEW`, `PENDING_PAYMENT`, `PAID`, `PROCESSING`, `SHIPPING`, `COMPLETED`, `CANCELLED`
- `payment_status`: `UNPAID`, `PAID`, `FAILED`, `REFUNDED`
- `payment_txn_type`: `PAYMENT`, `REFUND`
- `payment_txn_state`: `PENDING`, `SUCCESS`, `FAILED`, `CANCELLED`
- `notification_status`: `PENDING`, `SENT`, `FAILED`

## 3. Chi tiết bảng

### 3.1 `users`
| Column | Type | Null | Default | Constraint/Index | Mô tả |
|---|---|---|---|---|---|
| id | BINARY(16) | No | - | PK | UUID người dùng |
| email | VARCHAR(255) | No | - | UK (`uk_users__email`) | Email đăng nhập duy nhất |
| phone | VARCHAR(20) | No | - | UK (`uk_users__phone`) | SĐT duy nhất |
| password_hash | VARCHAR(255) | No | - | - | Mật khẩu băm |
| full_name | VARCHAR(150) | No | - | - | Tên hiển thị |
| role | ENUM('USER','ADMIN') | No | USER | IDX (`idx_users__role_status`) | Vai trò hệ thống |
| status | ENUM('ACTIVE','LOCKED','DISABLED') | No | ACTIVE | IDX (`idx_users__role_status`) | Trạng thái tài khoản |
| failed_login_attempts | TINYINT UNSIGNED | No | 0 | - | Số lần đăng nhập sai liên tiếp |
| locked_until | DATETIME(3) | Yes | NULL | - | Khóa tạm đến thời điểm |
| last_login_at | DATETIME(3) | Yes | NULL | - | Lần đăng nhập cuối |
| created_at | DATETIME(3) | No | - | - | Thời điểm tạo |
| updated_at | DATETIME(3) | No | - | - | Thời điểm cập nhật |
| deleted_at | DATETIME(3) | Yes | NULL | IDX (`idx_users__deleted_at`) | Soft delete |

### 3.2 `user_addresses`
| Column | Type | Null | Default | Constraint/Index | Mô tả |
|---|---|---|---|---|---|
| id | BINARY(16) | No | - | PK | UUID địa chỉ |
| user_id | BINARY(16) | No | - | FK -> users.id, IDX (`idx_user_addresses__user_id`) | Chủ địa chỉ |
| recipient_name | VARCHAR(150) | No | - | - | Người nhận |
| recipient_phone | VARCHAR(20) | No | - | - | SĐT nhận hàng |
| line1 | VARCHAR(255) | No | - | - | Địa chỉ dòng 1 |
| line2 | VARCHAR(255) | Yes | NULL | - | Địa chỉ dòng 2 |
| ward | VARCHAR(120) | Yes | NULL | - | Phường/xã |
| district | VARCHAR(120) | Yes | NULL | - | Quận/huyện |
| city | VARCHAR(120) | No | - | - | Thành phố |
| province | VARCHAR(120) | Yes | NULL | - | Tỉnh |
| postal_code | VARCHAR(20) | Yes | NULL | - | Mã bưu chính |
| country_code | CHAR(2) | No | VN | - | Quốc gia |
| is_default | TINYINT(1) | No | 0 | IDX (`idx_user_addresses__user_default`) | Địa chỉ mặc định |
| created_at | DATETIME(3) | No | - | - | Thời điểm tạo |
| updated_at | DATETIME(3) | No | - | - | Thời điểm cập nhật |
| deleted_at | DATETIME(3) | Yes | NULL | - | Soft delete |

### 3.3 `products`
| Column | Type | Null | Default | Constraint/Index | Mô tả |
|---|---|---|---|---|---|
| id | BINARY(16) | No | - | PK | UUID sản phẩm |
| sku | VARCHAR(64) | No | - | UK (`uk_products__sku`) | Mã SKU duy nhất |
| name | VARCHAR(255) | No | - | IDX (`idx_products__name`) | Tên sản phẩm |
| description | TEXT | Yes | NULL | - | Mô tả |
| price | DECIMAL(18,2) | No | - | CHECK (`price > 0`) | Giá bán |
| stock | INT UNSIGNED | No | 0 | CHECK (`stock >= 0`) | Tồn kho |
| status | ENUM('ACTIVE','INACTIVE','DISCONTINUED') | No | ACTIVE | IDX (`idx_products__status`) | Trạng thái bán |
| image_url | VARCHAR(500) | Yes | NULL | - | Ảnh đại diện |
| created_at | DATETIME(3) | No | - | - | Thời điểm tạo |
| updated_at | DATETIME(3) | No | - | - | Thời điểm cập nhật |
| deleted_at | DATETIME(3) | Yes | NULL | - | Soft delete |

### 3.4 `orders`
| Column | Type | Null | Default | Constraint/Index | Mô tả |
|---|---|---|---|---|---|
| id | BINARY(16) | No | - | PK | UUID đơn hàng |
| order_code | VARCHAR(32) | No | - | UK (`uk_orders__order_code`) | Mã đơn duy nhất |
| user_id | BINARY(16) | No | - | FK -> users.id, IDX (`idx_orders__user_created`) | Người mua |
| order_status | ENUM('NEW','PENDING_PAYMENT','PAID','PROCESSING','SHIPPING','COMPLETED','CANCELLED') | No | NEW | IDX (`idx_orders__status_created`) | Trạng thái đơn |
| payment_status | ENUM('UNPAID','PAID','FAILED','REFUNDED') | No | UNPAID | IDX (`idx_orders__payment_created`) | Trạng thái thanh toán |
| currency_code | CHAR(3) | No | VND | - | Loại tiền |
| subtotal_amount | DECIMAL(18,2) | No | - | CHECK (`subtotal_amount >= 0`) | Tạm tính |
| discount_amount | DECIMAL(18,2) | No | 0 | CHECK (`discount_amount >= 0`) | Giảm giá |
| shipping_fee | DECIMAL(18,2) | No | 0 | CHECK (`shipping_fee >= 0`) | Phí vận chuyển |
| total_amount | DECIMAL(18,2) | No | - | CHECK (`total_amount >= 0`) | Tổng tiền thanh toán |
| payment_due_at | DATETIME(3) | Yes | NULL | IDX (`idx_orders__payment_due`) | Deadline giữ đơn |
| customer_note | VARCHAR(500) | Yes | NULL | - | Ghi chú khách hàng |
| cancel_reason | VARCHAR(255) | Yes | NULL | - | Lý do hủy |
| shipping_recipient_name | VARCHAR(150) | No | - | - | Snapshot tên nhận hàng |
| shipping_phone | VARCHAR(20) | No | - | - | Snapshot SĐT nhận hàng |
| shipping_line1 | VARCHAR(255) | No | - | - | Snapshot địa chỉ dòng 1 |
| shipping_line2 | VARCHAR(255) | Yes | NULL | - | Snapshot địa chỉ dòng 2 |
| shipping_ward | VARCHAR(120) | Yes | NULL | - | Snapshot phường/xã |
| shipping_district | VARCHAR(120) | Yes | NULL | - | Snapshot quận/huyện |
| shipping_city | VARCHAR(120) | No | - | - | Snapshot thành phố |
| shipping_province | VARCHAR(120) | Yes | NULL | - | Snapshot tỉnh |
| shipping_postal_code | VARCHAR(20) | Yes | NULL | - | Snapshot mã bưu chính |
| shipping_country_code | CHAR(2) | No | VN | - | Snapshot quốc gia |
| placed_at | DATETIME(3) | No | - | - | Thời điểm đặt đơn |
| created_at | DATETIME(3) | No | - | - | Thời điểm tạo |
| updated_at | DATETIME(3) | No | - | - | Thời điểm cập nhật |

### 3.5 `order_items`
| Column | Type | Null | Default | Constraint/Index | Mô tả |
|---|---|---|---|---|---|
| id | BINARY(16) | No | - | PK | UUID dòng đơn hàng |
| order_id | BINARY(16) | No | - | FK -> orders.id, IDX (`idx_order_items__order_id`) | Đơn hàng cha |
| product_id | BINARY(16) | No | - | FK -> products.id | Sản phẩm tham chiếu |
| product_sku_snapshot | VARCHAR(64) | No | - | - | Snapshot SKU |
| product_name_snapshot | VARCHAR(255) | No | - | - | Snapshot tên sản phẩm |
| unit_price | DECIMAL(18,2) | No | - | CHECK (`unit_price > 0`) | Đơn giá tại thời điểm mua |
| quantity | INT UNSIGNED | No | - | CHECK (`quantity > 0`) | Số lượng |
| line_total | DECIMAL(18,2) | No | - | CHECK (`line_total >= 0`) | Thành tiền dòng |
| created_at | DATETIME(3) | No | - | - | Thời điểm tạo |
| updated_at | DATETIME(3) | No | - | - | Thời điểm cập nhật |

Ràng buộc bổ sung: `UNIQUE(order_id, product_id)`.

### 3.6 `order_status_history`
| Column | Type | Null | Default | Constraint/Index | Mô tả |
|---|---|---|---|---|---|
| id | BINARY(16) | No | - | PK | UUID lịch sử trạng thái |
| order_id | BINARY(16) | No | - | FK -> orders.id, IDX (`idx_osh__order_time`) | Đơn hàng |
| from_status | ENUM('NEW','PENDING_PAYMENT','PAID','PROCESSING','SHIPPING','COMPLETED','CANCELLED') | Yes | NULL | - | Trạng thái trước |
| to_status | ENUM('NEW','PENDING_PAYMENT','PAID','PROCESSING','SHIPPING','COMPLETED','CANCELLED') | No | - | IDX (`idx_osh__to_status_time`) | Trạng thái sau |
| changed_by_type | ENUM('SYSTEM','USER','ADMIN','PAYMENT_GATEWAY','CARRIER') | No | - | - | Tác nhân thay đổi |
| changed_by_id | BINARY(16) | Yes | NULL | - | ID tác nhân (nếu có) |
| change_reason | VARCHAR(255) | Yes | NULL | - | Lý do thay đổi |
| occurred_at | DATETIME(3) | No | - | IDX (`idx_osh__order_time`) | Thời điểm nghiệp vụ |
| created_at | DATETIME(3) | No | - | - | Thời điểm ghi log |

### 3.7 `payment_transactions`
| Column | Type | Null | Default | Constraint/Index | Mô tả |
|---|---|---|---|---|---|
| id | BINARY(16) | No | - | PK | UUID giao dịch |
| order_id | BINARY(16) | No | - | FK -> orders.id, IDX (`idx_payment_txn__order_status`) | Đơn liên quan |
| txn_type | ENUM('PAYMENT','REFUND') | No | PAYMENT | - | Loại giao dịch |
| provider_name | VARCHAR(50) | No | - | - | Tên cổng thanh toán |
| provider_txn_code | VARCHAR(128) | No | - | UK (`uk_payment_txn__provider_txn_code`) | Mã giao dịch cổng |
| provider_order_ref | VARCHAR(128) | Yes | NULL | - | Mã đơn phía cổng |
| idempotency_key | VARCHAR(128) | Yes | NULL | UK (`uk_payment_txn__idempotency_key`) | Chống xử lý lặp |
| amount | DECIMAL(18,2) | No | - | CHECK (`amount > 0`) | Số tiền |
| currency_code | CHAR(3) | No | VND | - | Loại tiền |
| status | ENUM('PENDING','SUCCESS','FAILED','CANCELLED') | No | PENDING | IDX (`idx_payment_txn__order_status`) | Trạng thái xử lý |
| provider_response_code | VARCHAR(64) | Yes | NULL | - | Mã phản hồi cổng |
| provider_message | VARCHAR(500) | Yes | NULL | - | Thông điệp phản hồi |
| verified_signature | TINYINT(1) | No | 0 | - | Kết quả verify chữ ký callback |
| payload_json | JSON | Yes | NULL | - | Raw payload (nếu cần lưu) |
| processed_at | DATETIME(3) | Yes | NULL | - | Thời điểm xử lý xong |
| created_at | DATETIME(3) | No | - | IDX (`idx_payment_txn__created_at`) | Thời điểm tạo |
| updated_at | DATETIME(3) | No | - | - | Thời điểm cập nhật |

### 3.8 `shipment_tracking_events`
| Column | Type | Null | Default | Constraint/Index | Mô tả |
|---|---|---|---|---|---|
| id | BINARY(16) | No | - | PK | UUID sự kiện giao vận |
| order_id | BINARY(16) | No | - | FK -> orders.id, IDX (`idx_ship_evt__order_time`) | Đơn liên quan |
| tracking_no | VARCHAR(64) | Yes | NULL | IDX (`idx_ship_evt__tracking_time`) | Mã vận đơn |
| carrier_code | VARCHAR(32) | Yes | NULL | - | Mã hãng vận chuyển |
| event_status | ENUM('PICKED_UP','IN_TRANSIT','OUT_FOR_DELIVERY','DELIVERED','DELIVERY_FAILED','RETURNED') | No | - | - | Trạng thái giao vận |
| event_description | VARCHAR(500) | Yes | NULL | - | Mô tả sự kiện |
| event_time | DATETIME(3) | No | - | IDX (`idx_ship_evt__order_time`) | Thời điểm sự kiện |
| source_type | ENUM('SYSTEM','ADMIN','CARRIER_WEBHOOK') | No | - | - | Nguồn cập nhật |
| source_ref | VARCHAR(128) | Yes | NULL | - | Mã tham chiếu nguồn |
| created_at | DATETIME(3) | No | - | - | Thời điểm ghi nhận |

### 3.9 `notifications`
| Column | Type | Null | Default | Constraint/Index | Mô tả |
|---|---|---|---|---|---|
| id | BINARY(16) | No | - | PK | UUID thông báo |
| recipient_user_id | BINARY(16) | No | - | FK -> users.id | Người nhận |
| channel | ENUM('IN_APP','EMAIL','SMS') | No | IN_APP | - | Kênh gửi |
| event_type | VARCHAR(64) | No | - | - | Loại sự kiện |
| event_key | VARCHAR(128) | No | - | - | Khóa dedupe sự kiện |
| title | VARCHAR(255) | No | - | - | Tiêu đề |
| content | TEXT | No | - | - | Nội dung |
| status | ENUM('PENDING','SENT','FAILED') | No | PENDING | IDX (`idx_notifications__status_created`) | Trạng thái gửi |
| sent_at | DATETIME(3) | Yes | NULL | - | Thời điểm gửi thành công |
| failure_reason | VARCHAR(255) | Yes | NULL | - | Lý do thất bại |
| created_at | DATETIME(3) | No | - | IDX (`idx_notifications__status_created`) | Thời điểm tạo |
| updated_at | DATETIME(3) | No | - | - | Thời điểm cập nhật |

Ràng buộc dedupe: `UNIQUE(recipient_user_id, channel, event_type, event_key)`.

### 3.10 `audit_logs`
| Column | Type | Null | Default | Constraint/Index | Mô tả |
|---|---|---|---|---|---|
| id | BINARY(16) | No | - | PK | UUID log |
| entity_type | ENUM('ORDER','PAYMENT_TRANSACTION','PRODUCT','USER') | No | - | IDX (`idx_audit__entity_time`) | Loại thực thể |
| entity_id | BINARY(16) | No | - | IDX (`idx_audit__entity_time`) | ID thực thể |
| action | VARCHAR(64) | No | - | IDX (`idx_audit__action_time`) | Hành động |
| actor_type | ENUM('SYSTEM','USER','ADMIN','PAYMENT_GATEWAY','CARRIER') | No | - | IDX (`idx_audit__actor_time`) | Loại tác nhân |
| actor_id | BINARY(16) | Yes | NULL | IDX (`idx_audit__actor_time`) | ID tác nhân |
| before_data | JSON | Yes | NULL | - | Snapshot trước |
| after_data | JSON | Yes | NULL | - | Snapshot sau |
| metadata | JSON | Yes | NULL | - | Thông tin bổ sung |
| ip_address | VARCHAR(45) | Yes | NULL | - | IP nguồn |
| user_agent | VARCHAR(255) | Yes | NULL | - | User agent |
| created_at | DATETIME(3) | No | - | IDX (`idx_audit__entity_time`) | Thời điểm log |

## 4. Ràng buộc liên bảng bắt buộc
- `orders.user_id` phải tồn tại trong `users`.
- `order_items.order_id` và `order_items.product_id` phải hợp lệ.
- `payment_transactions.order_id` phải hợp lệ.

## 5. Ràng buộc nên xử lý ở service layer
- Rule chuyển trạng thái đơn theo luồng nghiệp vụ chuẩn.
- Kiểm tra role khi ghi audit thao tác admin.
- Đảm bảo chỉ 1 địa chỉ mặc định/user tại một thời điểm.
