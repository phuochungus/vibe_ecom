# Versioning Strategy - Golf Store

## 1. Chuẩn version
- Dùng `Semantic Versioning`: `MAJOR.MINOR.PATCH`.
- `MAJOR`: thay đổi breaking (ảnh hưởng API/luồng nghiệp vụ cũ).
- `MINOR`: thêm feature mới tương thích ngược.
- `PATCH`: sửa lỗi, tối ưu nhỏ, không đổi hành vi chính.

## 2. Quy tắc gán version theo Agile
- Mỗi sprint có 1 release chính (ít nhất `MINOR` hoặc `MAJOR` nếu là MVP lần đầu).
- Trong sprint:
1. Build nội bộ: `vX.Y.Z-alpha.n`.
2. Build kiểm thử/UAT: `vX.Y.Z-rc.n`.
3. Release production: `vX.Y.Z`.

## 3. Version hiện tại
- Baseline trước sprint: `v0.9.0` (pre-MVP).
- Mục tiêu Sprint-01: phát hành `v1.0.0` vì hoàn thiện full scope MVP đã cam kết.

## 4. Quy tắc cập nhật tài liệu
- `BRD.md`: tăng version khi thay đổi scope/yêu cầu/AC.
- `business.md`: cập nhật khi thay đổi business rule.
- `PRODUCT_BACKLOG.md` và `SPRINT_01_PLAN.md`: cập nhật theo iteration, không cần SemVer riêng.
- Mọi thay đổi phát hành phải ghi vào `CHANGELOG.md`.

## 5. Quy tắc đặt tag Git
- Format tag: `release/vX.Y.Z`.
- Ví dụ Sprint-01: `release/v1.0.0`.
- Tag chỉ tạo sau khi pass UAT và có sign-off PO.
