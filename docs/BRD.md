# BRD - Hệ thống cửa hàng gậy golf (Web)

## 1. Thông tin tài liệu
- Mã tài liệu: `BRD-GOLF-STORE-001`
- Phiên bản: `1.1`
- Ngày ban hành: `26/02/2026`
- Người soạn: `BA/PO`
- Trạng thái: `Baseline for Sprint-01`

## 2. Mục tiêu kinh doanh
- Số hóa quy trình bán gậy golf trực tuyến cho khách hàng.
- Tăng tỷ lệ chốt đơn nhờ thanh toán tự động và theo dõi đơn hàng minh bạch.
- Nâng cao vận hành nội bộ qua quản lý sản phẩm, đơn hàng, khiếu nại và báo cáo doanh thu tập trung.

## 3. KPI nghiệp vụ
- Tỷ lệ thanh toán thành công >= 90%.
- Tỷ lệ đơn hoàn tất/đơn tạo mới >= 80%.
- Tỷ lệ giao hàng đúng hạn >= 95%.
- 95% khiếu nại được phản hồi lần đầu trong 48 giờ làm việc.

## 4. Phạm vi
### 4.1 In Scope
- Đăng nhập/đăng xuất user.
- Quản lý sản phẩm (admin).
- Tạo và quản lý đơn hàng.
- Thanh toán tự động qua cổng thanh toán.
- Tracking trạng thái đơn hàng.
- Thông báo theo sự kiện nghiệp vụ.
- Báo cáo doanh thu theo kỳ.
- Quản lý khiếu nại.

### 4.2 Out of Scope
- Quản lý chương trình loyalty/điểm thưởng.
- Tích hợp đa kho nâng cao (WMS).
- CSKH qua tổng đài tích hợp VOIP.

## 5. Đối tượng sử dụng và stakeholder
- `User`: khách mua hàng, theo dõi đơn và khiếu nại.
- `Admin`: quản lý sản phẩm, đơn hàng, khiếu nại, báo cáo.
- `Kế toán` (stakeholder gián tiếp): đối soát thanh toán và doanh thu.
- `Vận hành` (stakeholder gián tiếp): theo dõi xử lý đơn và SLA.

## 6. Giả định và phụ thuộc
- Có sẵn cổng thanh toán hỗ trợ callback/webhook.
- Có đơn vị vận chuyển hoặc trạng thái giao hàng được cập nhật vào hệ thống.
- User có email hoặc số điện thoại hợp lệ để nhận thông báo.

## 7. Quy trình nghiệp vụ tổng quan
### 7.1 Quy trình mua hàng
1. User đăng nhập.
2. User chọn sản phẩm và tạo đơn.
3. Hệ thống tạo đơn `Chờ thanh toán`.
4. User thanh toán online.
5. Hệ thống nhận kết quả thanh toán tự động.
6. Admin xử lý và giao hàng.
7. User theo dõi đơn đến khi `Hoàn tất`.

### 7.2 Quy trình khiếu nại
1. User tạo khiếu nại từ đơn hàng.
2. Admin tiếp nhận và phân loại.
3. Admin xử lý và phản hồi.
4. User xác nhận/đồng thuận kết quả.
5. Khiếu nại được đóng.

## 8. Danh sách Use Case
### UC-01: User đăng nhập
- Actor: `User`
- Tiền điều kiện: Tài khoản đã tồn tại.
- Luồng chính:
1. User nhập thông tin đăng nhập.
2. Hệ thống xác thực.
3. Hệ thống tạo phiên đăng nhập thành công.
- Luồng ngoại lệ:
1. Sai thông tin đăng nhập.
2. Tài khoản bị khóa tạm thời do đăng nhập sai nhiều lần.
- Hậu điều kiện: User đăng nhập thành công và có quyền mua hàng.

### UC-02: Admin quản lý sản phẩm
- Actor: `Admin`
- Tiền điều kiện: Admin đã đăng nhập.
- Luồng chính:
1. Tạo mới/sửa thông tin sản phẩm.
2. Cập nhật giá, tồn kho, trạng thái bán.
3. Lưu thay đổi.
- Luồng ngoại lệ:
1. Thiếu trường bắt buộc.
2. Giá hoặc tồn kho không hợp lệ.
- Hậu điều kiện: Danh mục sản phẩm được cập nhật.

### UC-03: User tạo đơn và thanh toán tự động
- Actor: `User`
- Tiền điều kiện: User đăng nhập, giỏ hàng có sản phẩm hợp lệ.
- Luồng chính:
1. User xác nhận checkout.
2. Hệ thống tạo mã đơn duy nhất.
3. User chọn phương thức thanh toán online.
4. Cổng thanh toán trả kết quả.
5. Hệ thống cập nhật trạng thái thanh toán tự động.
- Luồng ngoại lệ:
1. Thanh toán thất bại.
2. Quá thời gian giữ đơn chưa thanh toán.
- Hậu điều kiện: Đơn vào luồng xử lý giao vận hoặc hủy theo quy tắc.

### UC-04: User tracking đơn hàng
- Actor: `User`
- Tiền điều kiện: Có đơn hàng hợp lệ.
- Luồng chính:
1. User mở chi tiết đơn.
2. Hệ thống hiển thị trạng thái hiện tại + lịch sử trạng thái.
- Hậu điều kiện: User nắm được tiến độ đơn hàng.

### UC-05: Admin quản lý đơn hàng
- Actor: `Admin`
- Tiền điều kiện: Đơn đã tạo.
- Luồng chính:
1. Admin xác nhận đơn.
2. Admin cập nhật trạng thái xử lý/giao hàng.
3. Admin hoàn tất/hủy/hoàn theo chính sách.
- Hậu điều kiện: Đơn hàng được xử lý theo SLA.

### UC-06: Hệ thống gửi thông báo
- Actor: `System`
- Sự kiện kích hoạt: đặt hàng, thanh toán, đổi trạng thái đơn, cập nhật khiếu nại.
- Kết quả: User nhận thông báo đúng sự kiện, đúng người nhận.

### UC-07: Admin xem báo cáo doanh thu
- Actor: `Admin`
- Tiền điều kiện: Có dữ liệu đơn hàng.
- Luồng chính:
1. Chọn kỳ báo cáo.
2. Lọc theo trạng thái đơn.
3. Hệ thống trả dữ liệu tổng hợp + chi tiết theo đơn.
- Hậu điều kiện: Admin có dữ liệu theo dõi doanh thu.

### UC-08: User/Admin quản lý khiếu nại
- Actor: `User`, `Admin`
- Tiền điều kiện: User có đơn hàng hợp lệ.
- Luồng chính:
1. User tạo khiếu nại.
2. Admin tiếp nhận, xử lý, phản hồi.
3. Hệ thống cập nhật trạng thái khiếu nại.
4. Đóng khiếu nại.
- Hậu điều kiện: Vụ việc được xử lý và lưu vết.

## 9. Yêu cầu chức năng chi tiết (Functional Requirements)
- `FR-ACC-01`: Hệ thống cho phép user đăng nhập/đăng xuất.
- `FR-ACC-02`: Chỉ user đăng nhập mới được checkout, theo dõi đơn, gửi khiếu nại.
- `FR-PROD-01`: Admin thêm/sửa/xóa sản phẩm.
- `FR-PROD-02`: Admin cập nhật giá, tồn kho, trạng thái bán.
- `FR-ORD-01`: Hệ thống tạo mã đơn hàng duy nhất.
- `FR-ORD-02`: Hệ thống xử lý luồng trạng thái đơn hàng chuẩn.
- `FR-PAY-01`: Hệ thống nhận callback/webhook và cập nhật thanh toán tự động.
- `FR-PAY-02`: Hỗ trợ timeout giữ đơn và tự động hủy khi chưa thanh toán.
- `FR-TRK-01`: Hiển thị timeline trạng thái đơn theo thời gian.
- `FR-NOTI-01`: Tự động gửi thông báo theo sự kiện nghiệp vụ.
- `FR-REV-01`: Admin xem báo cáo doanh thu theo ngày/tháng/năm.
- `FR-REV-02`: Báo cáo cho phép lọc theo trạng thái đơn và truy vết theo mã đơn.
- `FR-CMP-01`: User tạo khiếu nại gắn với đơn hàng.
- `FR-CMP-02`: Admin xử lý, phản hồi, đóng khiếu nại theo SLA.

## 10. Business Rules
- `BR-ACC-01`: Email/số điện thoại user là duy nhất.
- `BR-ACC-02`: User chưa đăng nhập không được checkout.
- `BR-ACC-03`: Khóa tạm tài khoản 15 phút sau 5 lần đăng nhập sai liên tiếp.
- `BR-PROD-01`: Giá bán > 0.
- `BR-PROD-02`: Tồn kho >= 0.
- `BR-PROD-03`: Sản phẩm `Ngừng bán` không được thêm mới vào giỏ.
- `BR-ORD-01`: Mã đơn hàng là duy nhất toàn hệ thống.
- `BR-ORD-02`: Đơn hợp lệ phải có ít nhất 1 sản phẩm hợp lệ.
- `BR-ORD-03`: Trạng thái đơn chuẩn: `Mới tạo` -> `Chờ thanh toán` -> `Đã thanh toán` -> `Đang xử lý` -> `Đang giao` -> `Hoàn tất`.
- `BR-ORD-04`: Chỉ hủy đơn trước `Đang giao`, trừ xử lý ngoại lệ bởi admin.
- `BR-PAY-01`: Chỉ ghi nhận thanh toán khi callback hợp lệ.
- `BR-PAY-02`: Thanh toán thất bại chuyển `Chờ thanh toán` hoặc `Hủy` theo timeout.
- `BR-PAY-03`: Quá 30 phút chưa thanh toán thì tự động `Hủy`.
- `BR-PAY-04`: Mọi giao dịch thanh toán bắt buộc có mã giao dịch.
- `BR-TRK-01`: Mọi thay đổi trạng thái đơn phải lưu log trạng thái, thời điểm, tác nhân.
- `BR-NOTI-01`: Thông báo gửi trong vòng tối đa 1 phút từ khi phát sinh sự kiện.
- `BR-NOTI-02`: Không gửi trùng thông báo cùng sự kiện/cùng user.
- `BR-REV-01`: Doanh thu thuần chỉ tính đơn `Hoàn tất`.
- `BR-REV-02`: Đơn hoàn tiền bị trừ doanh thu trong kỳ.
- `BR-REV-03`: Báo cáo doanh thu phải truy vết tới mã đơn.
- `BR-CMP-01`: Khiếu nại phải gắn với đơn hợp lệ.
- `BR-CMP-02`: SLA xử lý khiếu nại mặc định 48 giờ làm việc.
- `BR-CMP-03`: Chỉ đóng khiếu nại sau phản hồi cuối cùng của admin.

## 11. Yêu cầu phi chức năng (NFR)
- `NFR-PERF-01`: API nghiệp vụ chính phản hồi < 3 giây ở tải thông thường.
- `NFR-SEC-01`: Mật khẩu lưu dạng băm; không lưu thông tin nhạy cảm thanh toán trái quy định.
- `NFR-AUD-01`: Có audit log cho thao tác đổi trạng thái đơn, thanh toán, khiếu nại.
- `NFR-AVL-01`: Có backup định kỳ và khả năng khôi phục dữ liệu.
- `NFR-REL-01`: Tỷ lệ thành công xử lý callback thanh toán >= 99%.

## 12. Mô hình dữ liệu nghiệp vụ mức cao
- `User`: thông tin định danh, trạng thái tài khoản.
- `Product`: mã, tên, giá, tồn kho, trạng thái bán.
- `Order`: mã đơn, user, tổng tiền, trạng thái đơn, trạng thái thanh toán.
- `PaymentTransaction`: mã giao dịch, mã đơn, số tiền, trạng thái, thời gian.
- `ShipmentTracking`: mã đơn, trạng thái, timestamp, nguồn cập nhật.
- `Complaint`: mã khiếu nại, mã đơn, nội dung, trạng thái, SLA.
- `Notification`: người nhận, loại sự kiện, trạng thái gửi, thời gian gửi.

## 13. Báo cáo nghiệp vụ
- `RPT-01`: Doanh thu theo ngày/tháng/năm.
- `RPT-02`: Tỷ lệ thanh toán thành công/thất bại.
- `RPT-03`: Tỷ lệ đơn hoàn tất, hủy, hoàn tiền.
- `RPT-04`: Tồn đọng khiếu nại theo SLA.

## 14. Acceptance Criteria tổng thể (UAT)
- `UAT-01`: User đăng nhập thành công mới thực hiện checkout.
- `UAT-02`: Đơn hàng được cập nhật trạng thái đúng luồng nghiệp vụ.
- `UAT-03`: Thanh toán online tự động cập nhật trạng thái trong hệ thống.
- `UAT-04`: User theo dõi được lịch sử đơn hàng theo timeline.
- `UAT-05`: Hệ thống gửi thông báo đúng sự kiện và không trùng.
- `UAT-06`: Admin quản lý được sản phẩm/đơn/khiếu nại từ giao diện quản trị.
- `UAT-07`: Báo cáo doanh thu phản ánh đúng dữ liệu đơn hoàn tất và hoàn tiền.
- `UAT-08`: Audit log truy vết được toàn bộ thay đổi trọng yếu.

## 15. Rủi ro và kiểm soát
- Rủi ro trễ callback thanh toán -> cơ chế retry + đối soát cuối ngày.
- Rủi ro sai lệch doanh thu do hoàn tiền -> rule hạch toán theo kỳ rõ ràng.
- Rủi ro quá tải khi flash sale -> giới hạn hàng đợi xử lý và giám sát hiệu năng.

## 16. Truy vết yêu cầu (RTM rút gọn)
- `UC-01` -> `FR-ACC-*` -> `BR-ACC-*` -> `UAT-01`.
- `UC-02` -> `FR-PROD-*` -> `BR-PROD-*` -> `UAT-06`.
- `UC-03` -> `FR-ORD-*`, `FR-PAY-*` -> `BR-ORD-*`, `BR-PAY-*` -> `UAT-02`, `UAT-03`.
- `UC-04` -> `FR-TRK-*` -> `BR-TRK-*` -> `UAT-04`.
- `UC-06` -> `FR-NOTI-*` -> `BR-NOTI-*` -> `UAT-05`.
- `UC-07` -> `FR-REV-*` -> `BR-REV-*` -> `UAT-07`.
- `UC-08` -> `FR-CMP-*` -> `BR-CMP-*` -> `UAT-06`.


## 17. Agile Delivery Links
- Product backlog: PRODUCT_BACKLOG.md.
- Sprint plan: SPRINT_01_PLAN.md.
- Versioning: VERSIONING.md.
- Release history: CHANGELOG.md.

