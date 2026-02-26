# Sprint-01 Plan - Golf Store

## 1. Sprint Information
- Sprint ID: `SPRINT-01`
- Sprint Goal: Hoàn thiện MVP end-to-end cho luồng bán hàng golf store (đăng nhập -> đặt hàng -> thanh toán -> giao vận -> tracking -> khiếu nại -> báo cáo).
- Thời gian sprint: `26/02/2026 - 11/03/2026` (2 tuần).
- Release mục tiêu: `v1.0.0`.

## 2. Team & Capacity
- Mô hình: Scrum.
- Vai trò:
1. `PO`: ưu tiên backlog, nghiệm thu nghiệp vụ.
2. `BA`: làm rõ requirement/AC/traceability.
3. `Dev`: thiết kế, code, unit test.
4. `QA`: test tích hợp, regression, UAT support.
- Capacity mục tiêu: `70 SP`.
- Planned scope: `68 SP`.

## 3. Sprint Scope (Committed)
- `US-01` Đăng nhập user.
- `US-02` Quản lý sản phẩm admin.
- `US-03` Quản lý đơn hàng.
- `US-04` Thanh toán tự động.
- `US-05` Tracking đơn hàng.
- `US-06` Thông báo hệ thống.
- `US-07` Báo cáo doanh thu.
- `US-08` Quản lý khiếu nại.

## 4. Sprint Board Workflow
- Trạng thái công việc:
1. `To Do`
2. `In Progress`
3. `Code Review`
4. `QA Testing`
5. `UAT`
6. `Done`
- WIP limit đề xuất:
1. `In Progress`: tối đa 3 ticket/dev.
2. `Code Review`: tối đa 5 ticket.
3. `QA Testing`: tối đa 6 ticket.

## 5. Ceremonies (lịch cố định)
- Sprint Planning: `26/02/2026`.
- Daily Scrum: mỗi ngày 15 phút.
- Backlog Refinement: `04/03/2026`.
- Sprint Review: `11/03/2026`.
- Retrospective: `11/03/2026`.

## 6. Definition of Ready (DoR)
- User story có mô tả rõ actor, mục tiêu, giá trị.
- Có acceptance criteria testable.
- Có mapping tới `UC/FR/BR` trong BRD.
- Không còn blocker nghiệp vụ mở.

## 7. Definition of Done (DoD)
- Code hoàn thành và pass code review.
- Unit/integration test pass.
- QA pass test case chính.
- UAT pass theo tiêu chí `UAT-01..UAT-08`.
- Có audit log cho thay đổi trạng thái đơn, thanh toán, khiếu nại.
- Tài liệu cập nhật: BRD, backlog, changelog.

## 8. Milestone kiểm soát
- `M1 - 01/03/2026`: xong US-01, US-02.
- `M2 - 05/03/2026`: xong US-03, US-04.
- `M3 - 08/03/2026`: xong US-05, US-06.
- `M4 - 10/03/2026`: xong US-07, US-08 + regression.
- `Release Candidate - 11/03/2026`: UAT sign-off và đóng sprint.

## 9. Quản lý rủi ro sprint
- Rủi ro tích hợp cổng thanh toán chậm callback.
- Rủi ro chồng chéo logic trạng thái đơn.
- Rủi ro mismatch dữ liệu báo cáo doanh thu.

## 10. Kế hoạch giảm rủi ro
- Thiết lập retry + idempotency cho callback thanh toán.
- Khóa rule chuyển trạng thái đơn ở service layer.
- Đối soát báo cáo theo mã đơn và transaction id hằng ngày.
