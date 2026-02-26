# Product Backlog - Golf Store (Agile)

## 1. Quy ước
- Estimation: Fibonacci Story Point (`3, 5, 8, 13`).
- Priority: `Must`, `Should`, `Could`.
- Sprint hiện tại: `Sprint-01`.

## 2. Epic List
- `EP-01`: Quản lý tài khoản và đăng nhập user.
- `EP-02`: Quản lý sản phẩm (admin).
- `EP-03`: Quản lý đơn hàng.
- `EP-04`: Thanh toán tự động.
- `EP-05`: Tracking đơn hàng.
- `EP-06`: Thông báo hệ thống.
- `EP-07`: Báo cáo doanh thu.

## 3. Sprint Backlog (Sprint-01)
| ID | Epic | User Story | Priority | SP | Sprint | Mapping BRD |
|---|---|---|---|---:|---|---|
| US-01 | EP-01 | Là user, tôi muốn đăng nhập để có thể mua hàng và theo dõi đơn. | Must | 5 | Sprint-01 | UC-01, FR-ACC-01/02, BR-ACC-02/03 |
| US-02 | EP-02 | Là admin, tôi muốn quản lý sản phẩm để luôn có dữ liệu bán hàng chính xác. | Must | 8 | Sprint-01 | UC-02, FR-PROD-01/02, BR-PROD-01/02/03 |
| US-03 | EP-03 | Là user/admin, tôi muốn tạo và xử lý đơn hàng theo trạng thái chuẩn. | Must | 13 | Sprint-01 | UC-03/05, FR-ORD-01/02, BR-ORD-01..04 |
| US-04 | EP-04 | Là user, tôi muốn thanh toán online và được cập nhật tự động kết quả. | Must | 13 | Sprint-01 | UC-03, FR-PAY-01/02, BR-PAY-01..04 |
| US-05 | EP-05 | Là user, tôi muốn xem tracking để biết đơn đang ở giai đoạn nào. | Must | 5 | Sprint-01 | UC-04, FR-TRK-01, BR-TRK-01 |
| US-06 | EP-06 | Là user, tôi muốn nhận thông báo khi có thay đổi đơn/thanh toán. | Must | 8 | Sprint-01 | UC-06, FR-NOTI-01, BR-NOTI-01/02 |
| US-07 | EP-07 | Là admin, tôi muốn xem báo cáo doanh thu theo kỳ để quản lý vận hành. | Must | 8 | Sprint-01 | UC-07, FR-REV-01/02, BR-REV-01..03 |

## 4. Acceptance Criteria tóm tắt theo User Story
### US-01
- Đăng nhập đúng thì tạo session thành công.
- Sai mật khẩu 5 lần liên tiếp thì khóa tạm 15 phút.

### US-02
- Admin tạo/sửa/xóa sản phẩm với trường bắt buộc hợp lệ.
- Không cho phép giá <= 0 hoặc tồn kho âm.
- Sản phẩm ngừng bán không xuất hiện để mua mới.

### US-03
- Tạo được mã đơn duy nhất.
- Đơn phải đi theo luồng trạng thái chuẩn.
- Hủy đơn chỉ trước trạng thái `Đang giao` trừ ngoại lệ admin.

### US-04
- Nhận callback hợp lệ thì cập nhật trạng thái thanh toán.
- Quá 30 phút chưa thanh toán thì tự hủy theo rule.
- Giao dịch luôn có mã để đối soát.

### US-05
- User xem được trạng thái hiện tại và timeline thay đổi.
- Mỗi thay đổi lưu trạng thái, thời điểm, tác nhân.

### US-06
- Thông báo gửi đúng sự kiện đơn hàng/thanh toán và đúng user.
- Không gửi trùng cho cùng sự kiện.

### US-07
- Báo cáo lọc theo ngày/tháng/năm.
- Doanh thu thuần chỉ tính đơn hoàn tất và trừ hoàn tiền.

## 5. Tổng effort Sprint-01
- Tổng story point committed: `60 SP`.
- Khuyến nghị năng lực team: tối thiểu `55-65 SP/sprint` cho scope hiện tại.
