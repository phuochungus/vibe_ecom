# Regression Test Cases - Golf Store

## 1. Scope
- Mục tiêu: dùng cho regression testing sau mỗi thay đổi backend/frontend, trước RC hoặc release.
- Phạm vi hiện tại: Sprint-01 MVP gồm `US-01` đến `US-07`.
- Hệ thống under test:
  - `FE_customer`
  - `FE_admin`
  - `BE_mono`

## 2. Environment
- Customer FE: `http://localhost:3001`
- Admin FE: `http://localhost:5173`
- Backend API: `http://localhost:8080`
- Database: PostgreSQL local

## 3. Test Accounts
- Admin:
  - Email: `admin@golf.local`
  - Password: `123456`
- User:
  - Email: `user@golf.local`
  - Password: `123456`

## 4. Execution Rule
- Mỗi test case phải ghi rõ:
  - `Pass` / `Fail`
  - ngày test
  - môi trường test
  - bug/ticket tham chiếu nếu fail
- Khi có thay đổi auth/order/payment/reporting, phải chạy full regression.

## 5. Regression Test Cases

| TC ID | Module | Test Case | Preconditions | Steps | Expected Result | Priority |
|---|---|---|---|---|---|---|
| REG-ACC-01 | Auth | User login thành công | User account tồn tại, active | 1. Mở FE customer. 2. Vào `/login`. 3. Nhập email/password hợp lệ. 4. Submit. | Đăng nhập thành công, token được lưu, chuyển về trang chủ, header hiển thị user. | High |
| REG-ACC-02 | Auth | Admin login thành công | Admin account tồn tại, active | 1. Mở FE admin. 2. Vào `/login`. 3. Nhập email/password hợp lệ. 4. Submit. | Đăng nhập thành công, vào dashboard admin. | High |
| REG-ACC-03 | Auth | Sai mật khẩu bị từ chối | Account tồn tại | 1. Vào login. 2. Nhập đúng email, sai password. 3. Submit. | API trả `401 INVALID_CREDENTIALS`, UI hiển thị lỗi phù hợp, không tạo session. | High |
| REG-ACC-04 | Auth | Account bị khóa tạm sau 5 lần sai mật khẩu | Account active | 1. Gửi 5 lần login sai liên tiếp. 2. Thử login lại. | Account bị khóa tạm 15 phút, API trả `423 ACCOUNT_LOCKED`. | High |
| REG-ACC-05 | Auth | Refresh token hoạt động | User/Admin đã login | 1. Login thành công. 2. Gọi refresh flow. | Access token mới được cấp, session vẫn hợp lệ. | Medium |
| REG-PROD-01 | Product | Admin xem danh sách sản phẩm | Admin đã login | 1. Vào menu `Sản phẩm`. | Danh sách hiển thị đúng dữ liệu, có SKU, tên, giá, tồn kho, trạng thái. | High |
| REG-PROD-02 | Product | Admin tạo sản phẩm hợp lệ | Admin đã login | 1. Vào `Sản phẩm`. 2. Chọn `Thêm sản phẩm`. 3. Nhập dữ liệu hợp lệ. 4. Submit. | Tạo thành công, quay về list, sản phẩm mới xuất hiện. | High |
| REG-PROD-03 | Product | Admin sửa sản phẩm | Admin đã login, có product tồn tại | 1. Mở chi tiết/chỉnh sửa sản phẩm. 2. Đổi giá hoặc tồn kho. 3. Submit. | Dữ liệu cập nhật đúng ở list/detail. | High |
| REG-PROD-04 | Product | Admin xóa sản phẩm | Admin đã login, có product tồn tại | 1. Chọn xóa sản phẩm. 2. Confirm. | Xóa thành công theo cơ chế hệ thống, sản phẩm không còn hiển thị trong list mua của user nếu trạng thái không cho bán. | High |
| REG-PROD-05 | Product | Validation giá không hợp lệ | Admin đã login | 1. Tạo/sửa sản phẩm với `price <= 0`. | Hệ thống chặn submit hoặc API trả validation error. | High |
| REG-PROD-06 | Product | Validation tồn kho âm | Admin đã login | 1. Tạo/sửa sản phẩm với `stock < 0`. | Hệ thống chặn submit hoặc API trả validation error. | High |
| REG-PROD-07 | Product | Product ngừng bán không xuất hiện cho user mua mới | Có product `DISCONTINUED` hoặc equivalent | 1. Set trạng thái ngừng bán. 2. Mở FE customer. 3. Kiểm tra list product. | Product không xuất hiện trong flow mua mới. | Medium |
| REG-ORD-01 | Order | User thêm sản phẩm vào giỏ | User đã login, có sản phẩm active | 1. Mở list sản phẩm. 2. Chọn add to cart. | Cart tăng số lượng, line item hiển thị đúng. | High |
| REG-ORD-02 | Order | User checkout COD thành công | User đã login, có item trong cart | 1. Vào cart. 2. Checkout. 3. Nhập địa chỉ hợp lệ. 4. Chọn COD. 5. Submit. | Đơn được tạo thành công, có `order_code`, chuyển đến trang kết quả hoặc orders. | High |
| REG-ORD-03 | Order | Order code là duy nhất | Có thể tạo đơn | 1. Tạo 2 đơn khác nhau. | Mỗi đơn có `order_code` khác nhau. | High |
| REG-ORD-04 | Order | User xem danh sách đơn hàng của mình | User đã login, có order | 1. Vào `/orders`. | Hiển thị đúng danh sách order của user, không lẫn user khác. | High |
| REG-ORD-05 | Order | User xem chi tiết đơn hàng | User đã login, có order | 1. Từ list order chọn `Xem chi tiết`. | Hiển thị đầy đủ items, amount, shipping info, payment state. | High |
| REG-ORD-06 | Order | User hủy đơn trước `SHIPPING` | User đã login, order đang ở trạng thái cho phép | 1. Mở order detail. 2. Thực hiện hủy. | Order chuyển `CANCELLED`, stock được hoàn lại theo rule. | High |
| REG-ORD-07 | Order | User không thể hủy đơn sau `SHIPPING` | Order đã ở `SHIPPING` hoặc cao hơn | 1. Thử hủy đơn. | API/UI chặn, trả lỗi `ORDER_CANNOT_CANCEL` hoặc tương đương. | High |
| REG-ORD-08 | Order | Admin xem danh sách đơn với filter status | Admin đã login, có dữ liệu nhiều trạng thái | 1. Vào `Đơn hàng`. 2. Filter theo `PAID`, `SHIPPING`, `COMPLETED`. | Kết quả khác nhau đúng theo trạng thái; status filter không bị ignore. | High |
| REG-ORD-09 | Order | Admin cập nhật trạng thái đơn đúng luồng | Admin đã login, order ở trạng thái phù hợp | 1. Mở order detail. 2. Chọn status tiếp theo hợp lệ. 3. Submit reason nếu cần. | Trạng thái cập nhật thành công, timeline được ghi nhận. | High |
| REG-ORD-10 | Order | Admin không thể cập nhật trạng thái sai luồng | Admin đã login | 1. Thử nhảy trạng thái không hợp lệ. | API trả conflict/validation error, dữ liệu không đổi. | High |
| REG-PAY-01 | Payment | Tạo payment online thành công | User đã login, có order unpaid | 1. Vào order detail hoặc checkout. 2. Chọn thanh toán online. | Hệ thống tạo payment transaction và trả `checkout_url`. | High |
| REG-PAY-02 | Payment | Callback thanh toán thành công cập nhật order | Có order unpaid và callback hợp lệ | 1. Gửi callback success. | `payment_status` và `order_status` cập nhật đúng theo rule. | High |
| REG-PAY-03 | Payment | Callback thanh toán thất bại cập nhật đúng | Có order unpaid | 1. Gửi callback failed. | Payment transaction được ghi nhận, trạng thái đơn/thanh toán cập nhật đúng. | High |
| REG-PAY-04 | Payment | Callback duplicate là idempotent | Có payment transaction đã xử lý | 1. Gửi lại cùng callback/transaction id. | Không double update, không tạo transaction trùng. | High |
| REG-PAY-05 | Payment | Đơn quá hạn chưa thanh toán bị expire theo rule | Có order pending payment quá hạn | 1. Trigger job/flow expire. | Đơn bị cancel/expire theo business rule, stock hoàn lại nếu áp dụng. | Medium |
| REG-TRK-01 | Tracking | User xem timeline đơn hàng | User đã login, có order | 1. Mở order detail. 2. Xem phần timeline/tracking. | Hiển thị trạng thái hiện tại, lịch sử thay đổi, thời điểm, tác nhân/description nếu có. | High |
| REG-TRK-02 | Tracking | Admin xem tracking timeline | Admin đã login | 1. Mở admin order detail. | Timeline hiển thị đúng với dữ liệu order. | Medium |
| REG-NOTI-01 | Notification | User xem danh sách thông báo | User đã login, có notification | 1. Vào `/notifications`. | Danh sách thông báo hiển thị đúng, không lỗi auth. | High |
| REG-NOTI-02 | Notification | Mark one notification as read | User đã login | 1. Mở notification. 2. Mark as read. | Notification đổi trạng thái read thành công. | Medium |
| REG-NOTI-03 | Notification | Mark all notifications as read | User đã login, có nhiều notification unread | 1. Chọn `Mark all as read`. | Tất cả notification unread chuyển sang read. | Medium |
| REG-NOTI-04 | Notification | Notification không bị gửi trùng cùng event | Có event order/payment lặp lại | 1. Trigger lại cùng event key. | Không tạo notification duplicate cho cùng user/event_key. | High |
| REG-REV-01 | Reporting | Admin xem revenue summary | Admin đã login | 1. Vào `Doanh thu`. | Summary hiển thị gross revenue, refund, net revenue, completed orders, payment success rate. | High |
| REG-REV-02 | Reporting | Revenue filter theo khoảng ngày | Admin đã login | 1. Chọn khoảng ngày khác nhau. | Dữ liệu summary/order list thay đổi tương ứng. | High |
| REG-REV-03 | Reporting | Net revenue tính đúng | Có order completed và refund data | 1. Đối soát dữ liệu summary với order/payment source. | `net_revenue = gross_revenue - refund_amount`. | High |
| REG-RBAC-01 | Security | User không truy cập được admin pages | User đã login | 1. Truy cập `FE_admin` hoặc admin API. | Bị chặn bởi UI/API. | High |
| REG-RBAC-02 | Security | Admin không dùng nhầm customer-only flow gây lỗi | Admin đã login | 1. Thử dùng customer protected endpoints không phù hợp. | Hệ thống xử lý đúng theo RBAC hoặc business rule. | Medium |
| REG-API-01 | Contract | Login API nhận field `email` | Backend đang chạy | 1. Gọi `POST /api/v1/auth/login` với `email/password`. | Trả đúng contract `docs/FE/API_CONTRACT.md`. | High |
| REG-API-02 | Contract | API envelope consistent | Backend đang chạy | 1. Gọi success endpoint. 2. Gọi fail endpoint. | Envelope success/error đúng format chuẩn. | Medium |
| REG-OPS-01 | Smoke | Health endpoint hoạt động | Backend đang chạy | 1. Gọi `/healthz`. 2. Gọi `/readyz`. | Trả healthy/ready đúng. | Medium |

## 6. Minimum Regression Suite Before Release
- Chạy bắt buộc:
  - `REG-ACC-01` đến `REG-ACC-04`
  - `REG-PROD-01` đến `REG-PROD-06`
  - `REG-ORD-02`, `REG-ORD-04`, `REG-ORD-05`, `REG-ORD-08`, `REG-ORD-09`
  - `REG-PAY-01` đến `REG-PAY-04`
  - `REG-TRK-01`
  - `REG-NOTI-01`, `REG-NOTI-03`, `REG-NOTI-04`
  - `REG-REV-01` đến `REG-REV-03`
  - `REG-RBAC-01`

## 7. Traceability
- `US-01` -> `REG-ACC-*`
- `US-02` -> `REG-PROD-*`
- `US-03` -> `REG-ORD-*`
- `US-04` -> `REG-PAY-*`
- `US-05` -> `REG-TRK-*`
- `US-06` -> `REG-NOTI-*`
- `US-07` -> `REG-REV-*`

## 8. Recommended Automation Mapping
- `@playwright/test`
  - login
  - admin product CRUD smoke
  - checkout
  - order detail/tracking
  - notifications
  - admin revenue
- API/integration test
  - duplicate callback idempotency
  - payment expiry job
  - revenue correctness
  - notification dedupe
  - RBAC negative cases
