# Reporting Model - Doanh thu và vận hành

## 1. Mục tiêu báo cáo
- Báo cáo doanh thu theo ngày/tháng/năm.
- Theo dõi tỷ lệ thanh toán thành công/thất bại.
- Theo dõi tỷ lệ đơn hoàn tất/hủy/hoàn tiền.

## 2. Nguồn dữ liệu chính
- `orders`: trạng thái đơn, tổng tiền, thời điểm đặt/hoàn tất.
- `payment_transactions`: thanh toán/hoàn tiền, trạng thái giao dịch.
- `order_status_history`: timeline thay đổi trạng thái (dùng để audit KPI vận hành).

## 3. Grain dữ liệu đề xuất
- Order fact grain: 1 dòng/đơn hàng (`orders.id`).
- Payment fact grain: 1 dòng/giao dịch (`payment_transactions.id`).

## 4. Công thức KPI chuẩn
### 4.1 Doanh thu thuần theo kỳ
- `Gross Revenue`: tổng `orders.total_amount` với `orders.order_status = COMPLETED`.
- `Refund Amount`: tổng `payment_transactions.amount` với `txn_type = REFUND` và `status = SUCCESS`.
- `Net Revenue`: `Gross Revenue - Refund Amount`.

### 4.2 Tỷ lệ thanh toán thành công
- Mẫu số: số giao dịch `txn_type = PAYMENT`.
- Tử số: số giao dịch `txn_type = PAYMENT AND status = SUCCESS`.
- `Success Rate = Tử số / Mẫu số`.

### 4.3 Tỷ lệ đơn hoàn tất
- Mẫu số: tổng số đơn phát sinh trong kỳ.
- Tử số: số đơn có `order_status = COMPLETED`.

## 5. Logic lọc thời gian
- Kỳ báo cáo mặc định dùng UTC trong DB.
- Nếu hiển thị theo timezone địa phương, conversion thực hiện ở service/report layer.
- Khuyến nghị khóa kỳ theo `[from_utc, to_utc)` để tránh trùng lặp biên.

## 6. Logical views đề xuất (không tạo SQL ở pha này)
1. `v_order_financial_fact`
- Cột chính: `order_id`, `order_code`, `order_status`, `payment_status`, `total_amount`, `placed_at`, `completed_at` (derive).

2. `v_payment_daily_summary`
- Cột chính: `report_date`, `payment_success_count`, `payment_failed_count`, `payment_success_amount`, `refund_amount`.

## 7. Quy tắc đối soát
- Mỗi bản ghi báo cáo doanh thu phải truy vết được về `orders.order_code`.
- Mọi refund ghi nhận trong kỳ phải tham chiếu `order_id` hợp lệ.
- Chênh lệch giữa payment success và order paid phải có báo cáo reconciliation.
