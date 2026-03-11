# Implementation Checklist - DB Design to Build

## 1. Chuẩn bị
- [ ] Xác nhận PostgreSQL version >= 16.
- [ ] Thiết lập encoding/timezone: `UTF8` / `UTC`.
- [ ] Chuẩn hóa UTC cho toàn bộ service kết nối DB.
- [ ] Thống nhất chiến lược UUID string giữa app và database.

## 2. Tạo schema (khi vào pha implement)
- [ ] Tạo bảng theo thứ tự phụ thuộc FK: `users` -> `products` -> `orders` -> `order_items` -> bảng log/transaction.
- [ ] Tạo đầy đủ PK/FK/UK/CHECK theo `DATA_DICTIONARY.md`.
- [ ] Tạo index phục vụ truy vấn chính (order list, tracking, reporting).
- [ ] Bật cơ chế transaction cho luồng tạo đơn và thanh toán callback.

## 3. Enforce business rules
- [ ] Tạo unique cho email/phone/order_code/provider_txn_code.
- [ ] Enforce check cho price/stock/quantity/amount.
- [ ] Tạo composite unique chống trùng notification.
- [ ] Implement state machine đơn hàng ở service layer.
- [ ] Implement scheduler timeout thanh toán (30 phút).

## 4. Test scenario bắt buộc
- [ ] Tạo user trùng email/phone -> fail.
- [ ] Tạo product `price <= 0` hoặc `stock < 0` -> fail.
- [ ] Tạo order không có item -> fail.
- [ ] Callback trùng `provider_txn_code` -> idempotent, không double update.
- [ ] Hủy đơn sau `SHIPPING` -> fail (trừ admin override có audit).
- [ ] Tracking trả đúng timeline từ `order_status_history` + `shipment_tracking_events`.
- [ ] Báo cáo net revenue = completed orders - successful refunds.
- [ ] Gửi notification trùng event_key -> bị dedupe.

## 5. NFR và vận hành
- [ ] Đảm bảo truy vấn chính có index phù hợp để đáp ứng SLA API < 3s.
- [ ] Ghi audit log cho thay đổi trạng thái đơn/thanh toán.
- [ ] Thiết lập chính sách backup/restore và test restore định kỳ.
- [ ] Thiết lập dashboard theo dõi callback success rate >= 99%.

## 6. Traceability
- [ ] Cross-check với `docs/BRD.md` các mục `FR-*`, `BR-*`, `UAT-*`.
- [ ] Cross-check với `docs/business.md` để không lệch business rule.
- [ ] Cập nhật `docs/CHANGELOG.md` khi schema chính thức được implement.
