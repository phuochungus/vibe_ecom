# Tài liệu yêu cầu nghiệp vụ (Business Requirements & Business Rules)

## 1. Mục tiêu hệ thống
Xây dựng hệ thống cửa hàng gậy golf cho phép khách hàng mua hàng trực tuyến và đội ngũ quản trị vận hành toàn bộ quy trình bán hàng, thanh toán, giao hàng và báo cáo doanh thu.

## 2. Phạm vi
- Kênh: Web app.
- Nhóm người dùng:
1. `User` (khách mua hàng).
2. `Admin` (quản lý vận hành hệ thống).

## 3. Vai trò và quyền hạn
### 3.1 User
- Đăng ký/đăng nhập.
- Xem sản phẩm, đặt hàng, thanh toán tự động.
- Theo dõi trạng thái đơn hàng.
- Nhận thông báo hệ thống.

### 3.2 Admin
- Quản lý sản phẩm (thêm/sửa/xóa, giá, tồn kho, trạng thái bán).
- Quản lý đơn hàng (xác nhận, xử lý, giao vận, hủy/hoàn).
- Theo dõi thanh toán và đối soát.
- Xem báo cáo doanh thu.
- Gửi thông báo đến user.

## 4. Yêu cầu nghiệp vụ chức năng
### 4.1 Đăng nhập người dùng
- Hệ thống cho phép user đăng ký, đăng nhập và đăng xuất.
- Chỉ user đã đăng nhập mới được đặt hàng và theo dõi đơn.

### 4.2 Quản lý sản phẩm (Admin)
- Admin tạo mới sản phẩm với thông tin bắt buộc: mã sản phẩm, tên, giá bán, tồn kho, trạng thái.
- Admin cập nhật thông tin sản phẩm theo thời gian thực.
- Sản phẩm ngừng kinh doanh không hiển thị cho user mới mua.

### 4.3 Quản lý đơn hàng
- User tạo đơn từ giỏ hàng.
- Hệ thống tạo mã đơn hàng duy nhất.
- Admin xử lý đơn theo luồng trạng thái chuẩn.

### 4.4 Thanh toán tự động
- Hệ thống tích hợp cổng thanh toán và tự động cập nhật kết quả thanh toán.
- Đơn hàng chỉ chuyển sang xử lý giao vận khi thanh toán thành công (đối với phương thức trả trước).

### 4.5 Tracking đơn hàng
- User xem được lịch sử và trạng thái hiện tại của đơn hàng.
- Mỗi lần thay đổi trạng thái đơn, hệ thống ghi nhận mốc thời gian và người/cơ chế cập nhật.

### 4.6 Thông báo
- Hệ thống gửi thông báo cho user khi: đặt hàng thành công, thanh toán thành công/thất bại, đơn chuyển trạng thái, đơn hủy/hoàn.

### 4.7 Báo cáo doanh thu
- Admin xem báo cáo doanh thu theo ngày/tháng/năm.
- Báo cáo có khả năng lọc theo khoảng thời gian và trạng thái đơn hàng.

## 5. Business Rules
### 5.1 Tài khoản và đăng nhập
- `BR-ACC-01`: Email/số điện thoại của user là duy nhất trong hệ thống.
- `BR-ACC-02`: User chưa đăng nhập không được phép checkout.
- `BR-ACC-03`: Sau 5 lần đăng nhập sai liên tiếp, tài khoản bị khóa tạm thời 15 phút.

### 5.2 Sản phẩm
- `BR-PROD-01`: Giá bán phải lớn hơn 0.
- `BR-PROD-02`: Tồn kho không được âm.
- `BR-PROD-03`: Sản phẩm ở trạng thái `Ngừng bán` không được thêm mới vào giỏ hàng.

### 5.3 Đơn hàng
- `BR-ORD-01`: Mỗi đơn hàng có mã duy nhất, không trùng lặp.
- `BR-ORD-02`: User chỉ được tạo đơn khi giỏ hàng có ít nhất 1 sản phẩm hợp lệ.
- `BR-ORD-03`: Luồng trạng thái đơn hàng chuẩn:
`Mới tạo` -> `Chờ thanh toán` -> `Đã thanh toán` -> `Đang xử lý` -> `Đang giao` -> `Hoàn tất`.
- `BR-ORD-04`: Đơn chỉ được hủy trước khi chuyển sang trạng thái `Đang giao` (trừ trường hợp admin xử lý ngoại lệ).

### 5.4 Thanh toán tự động
- `BR-PAY-01`: Kết quả thanh toán chỉ được ghi nhận khi nhận phản hồi hợp lệ từ cổng thanh toán.
- `BR-PAY-02`: Nếu thanh toán thất bại, đơn quay về `Chờ thanh toán` hoặc `Hủy` theo cấu hình timeout.
- `BR-PAY-03`: Nếu quá thời gian giữ đơn (ví dụ 30 phút) chưa thanh toán, đơn tự động chuyển `Hủy`.
- `BR-PAY-04`: Mọi giao dịch thanh toán phải có mã giao dịch để phục vụ đối soát.

### 5.5 Tracking và thông báo
- `BR-TRK-01`: Mỗi thay đổi trạng thái đơn phải lưu lịch sử (trạng thái, thời điểm, tác nhân).
- `BR-NOTI-01`: Thông báo gửi cho user tối đa trong vòng 1 phút sau khi phát sinh sự kiện nghiệp vụ.
- `BR-NOTI-02`: Không gửi trùng thông báo cho cùng một sự kiện và cùng người nhận.

### 5.6 Báo cáo doanh thu
- `BR-REV-01`: Doanh thu thuần chỉ tính từ đơn ở trạng thái `Hoàn tất`.
- `BR-REV-02`: Đơn hoàn tiền phải được trừ khỏi doanh thu trong kỳ tương ứng.
- `BR-REV-03`: Dữ liệu báo cáo phải truy vết được tới mã đơn hàng.

## 6. Phi chức năng (NFR)
- Hiệu năng: thời gian phản hồi API nghiệp vụ chính < 3 giây trong điều kiện tải thông thường.
- Bảo mật: mật khẩu lưu dưới dạng băm; dữ liệu thanh toán tuân thủ quy định bảo mật cổng thanh toán.
- Audit: toàn bộ thao tác thay đổi trạng thái đơn và thanh toán phải có nhật ký.
- Sẵn sàng: hệ thống hỗ trợ vận hành liên tục, có cơ chế backup dữ liệu định kỳ.

## 7. KPI đề xuất
- Tỷ lệ thanh toán thành công.
- Tỷ lệ giao hàng đúng hạn.
- Tỷ lệ đơn hoàn tất/đơn tạo mới.
- Doanh thu theo kỳ và tăng trưởng theo tháng.

## 8. Tiêu chí nghiệm thu mức nghiệp vụ
- User có thể đăng nhập, đặt hàng, thanh toán và theo dõi đơn hàng end-to-end.
- Admin quản lý được sản phẩm, đơn hàng và xem báo cáo doanh thu.
- Hệ thống gửi thông báo đúng sự kiện và đúng người nhận.
- Các business rule trọng yếu (`BR-ORD`, `BR-PAY`, `BR-REV`) được áp dụng nhất quán.
