# Admin Web (MVP)

## Mục tiêu
Trang quản trị nội bộ cho Golf Store để theo dõi:
- Tổng quan doanh thu và đơn hàng.
- Danh sách sản phẩm.
- Đơn hàng gần đây.
- Danh sách khách hàng nổi bật.
- Danh sách user admin và trạng thái hoạt động.

## Chạy thử
```bash
cd /workspace/golf_store
python3 -m http.server 4173
```
Sau đó mở: `http://localhost:4173/admin-web/`

## Chức năng đã có
- Điều hướng theo từng phân hệ: Tổng quan, Sản phẩm, Đơn hàng, Khách hàng, User Admin.
- Hiển thị KPI dashboard: doanh thu tháng, đơn mới, số admin đang hoạt động.
- Thêm sản phẩm nhanh từ form.
- Thêm user admin từ form (họ tên, email, vai trò).
- Khóa/mở khóa nhanh user admin bằng nút hành động.
- Tìm kiếm đồng thời trên sản phẩm, đơn hàng, khách hàng và user admin.

## Gợi ý bước tiếp theo
- Kết nối API thật theo `docs/FE/openapi.yaml`.
- Thêm đăng nhập và xác thực token.
- Bổ sung phân quyền chi tiết theo role.
- Bổ sung biểu đồ doanh thu theo ngày/tuần/tháng.
