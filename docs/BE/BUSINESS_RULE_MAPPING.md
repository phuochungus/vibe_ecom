# Business Rule Mapping - DB vs Service

## 1. Nguyên tắc phân lớp
- `DB layer`: đảm bảo toàn vẹn dữ liệu (PK/FK/UK/CHECK/index).
- `Service layer`: xử lý luồng nghiệp vụ, verify callback, authorization, scheduler.

## 2. Mapping chi tiết
| Rule ID | Mô tả ngắn | Cơ chế chính | Bảng/Cột liên quan | Ghi chú |
|---|---|---|---|---|
| BR-ACC-01 | Email/SĐT duy nhất | DB | `users.email`, `users.phone` (UNIQUE) | Chặn trùng ở tầng persistence |
| BR-ACC-02 | Chưa đăng nhập không checkout | Service | Auth middleware + order API | Không xử lý ở DB |
| BR-ACC-03 | Sai 5 lần khóa 15 phút | Service + DB field | `users.failed_login_attempts`, `users.locked_until` | Service cập nhật counter/lock |
| BR-PROD-01 | Giá > 0 | DB | `products.price` (CHECK) | Validate thêm ở API |
| BR-PROD-02 | Tồn kho không âm | DB | `products.stock` (UNSIGNED/CHECK) | Dùng transaction khi trừ kho |
| BR-PROD-03 | Ngừng bán không cho mua mới | Service | `products.status` | Check khi add cart/create order |
| BR-ORD-01 | Mã đơn duy nhất | DB | `orders.order_code` (UNIQUE) | Sinh mã ở service |
| BR-ORD-02 | Đơn phải có >=1 item hợp lệ | Service + DB | `order_items`, `orders` | Enforce trong transaction tạo đơn |
| BR-ORD-03 | Luồng trạng thái đơn chuẩn | Service + Audit | `orders.order_status`, `order_status_history` | Dùng state machine ở service |
| BR-ORD-04 | Hủy trước SHIPPING (trừ ngoại lệ admin) | Service + Audit | `orders.order_status`, `audit_logs` | Ghi audit khi override |
| BR-PAY-01 | Chỉ ghi nhận callback hợp lệ | Service | `payment_transactions.verified_signature` | Verify chữ ký trước commit |
| BR-PAY-02 | Thất bại -> chờ thanh toán/hủy | Service | `orders.payment_status`, `orders.order_status` | Theo timeout policy |
| BR-PAY-03 | Quá 30 phút tự hủy | Service job | `orders.payment_due_at` | Scheduler định kỳ |
| BR-PAY-04 | Bắt buộc có mã giao dịch | DB | `payment_transactions.provider_txn_code` | UNIQUE để idempotent |
| BR-TRK-01 | Lưu lịch sử đổi trạng thái | DB + Service | `order_status_history` | Mỗi transition phải insert log |
| BR-NOTI-01 | Thông báo trong 1 phút | Service + Queue | `notifications` | SLA vận hành, không enforce thuần DB |
| BR-NOTI-02 | Không gửi trùng cùng sự kiện | DB | UNIQUE(`recipient_user_id`,`channel`,`event_type`,`event_key`) | Dedupe cứng |
| BR-REV-01 | Doanh thu thuần từ đơn COMPLETED | Service/Reporting | `orders.order_status`, `orders.total_amount` | Dùng query/report view |
| BR-REV-02 | Hoàn tiền trừ doanh thu kỳ | Service/Reporting | `payment_transactions` (`txn_type=REFUND`,`status=SUCCESS`) | Hạch toán theo kỳ báo cáo |
| BR-REV-03 | Truy vết đến mã đơn | DB design | `orders.order_code`, FK liên bảng | Report phải giữ order-level drilldown |

## 3. Checklist ownership
- DB chịu trách nhiệm: integrity + uniqueness + anti-duplicate.
- Service chịu trách nhiệm: workflow, security, policy, scheduler.
- Reporting chịu trách nhiệm: công thức KPI, kỳ tính doanh thu, đối soát hoàn tiền.
