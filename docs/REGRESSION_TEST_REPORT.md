# Regression Test Report - Golf Store

## 1. Report Info
- Report date: `2026-03-11`
- Scope: regression verification theo [REGRESSION_TEST_CASES.md](/Users/imacvip/vibe_ecom/docs/REGRESSION_TEST_CASES.md)
- Environment:
  - `FE_admin`: `http://localhost:5173`
  - `FE_customer`: `http://localhost:3001`
  - `BE_mono`: `http://localhost:8080`
  - Database: PostgreSQL local
  - UI evidence directory: [tests](/Users/imacvip/vibe_ecom/tests)

## 2. Execution Summary
- Automated Playwright:
  - `FE_admin`: `4/4 passed`
  - `FE_customer`: `4/4 passed`
  - Video recording: `enabled`
- API smoke:
  - `GET /healthz`: `passed`
  - `POST /api/v1/auth/login`: `passed`

## 3. Test Case Result Matrix

| TC ID | Module | Status | Execution Type | Evidence / Notes |
|---|---|---|---|---|
| REG-ACC-01 | Auth | Pass | Automated | Video: [fe-customer--customer-can-log-in-and-browse-products--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-log-in-and-browse-products--passed.webm) |
| REG-ACC-02 | Auth | Pass | Automated | Video: [fe-admin--admin-can-log-in-and-see-dashboard-data--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-log-in-and-see-dashboard-data--passed.webm) |
| REG-ACC-03 | Auth | Not Run | Manual/API | Chưa chạy negative login sai mật khẩu. |
| REG-ACC-04 | Auth | Not Run | Manual/API | Chưa chạy lock sau 5 lần sai mật khẩu. |
| REG-ACC-05 | Auth | Not Run | Manual/API | Chưa chạy refresh token explicit regression. |
| REG-PROD-01 | Product | Pass | Automated | Covered by admin product navigation flow. Video: [fe-admin--admin-can-log-in-and-see-dashboard-data--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-log-in-and-see-dashboard-data--passed.webm) |
| REG-PROD-02 | Product | Pass | Automated | Video: [fe-admin--admin-can-create-a-product--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-create-a-product--passed.webm) |
| REG-PROD-03 | Product | Not Run | Manual/UI | Chưa chạy edit product regression. |
| REG-PROD-04 | Product | Not Run | Manual/UI | Chưa chạy delete product regression. |
| REG-PROD-05 | Product | Not Run | Manual/API | Chưa verify validation `price <= 0`. |
| REG-PROD-06 | Product | Not Run | Manual/API | Chưa verify validation `stock < 0`. |
| REG-PROD-07 | Product | Not Run | Manual/UI | Chưa verify product discontinued không hiển thị cho user mua mới. |
| REG-ORD-01 | Order | Pass | Automated | Covered by customer browse/add-to-cart flow. Video: [fe-customer--customer-can-log-in-and-browse-products--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-log-in-and-browse-products--passed.webm) |
| REG-ORD-02 | Order | Pass | Automated | Video: [fe-customer--customer-can-complete-cod-checkout--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-complete-cod-checkout--passed.webm) |
| REG-ORD-03 | Order | Not Run | Manual/API | Chưa đối soát uniqueness của `order_code`. |
| REG-ORD-04 | Order | Pass | Automated | Video: [fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm) |
| REG-ORD-05 | Order | Pass | Automated | Video: [fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm) |
| REG-ORD-06 | Order | Not Run | Manual/API | Chưa chạy user cancel order happy path. |
| REG-ORD-07 | Order | Not Run | Manual/API | Chưa chạy cancel blocked after `SHIPPING`. |
| REG-ORD-08 | Order | Not Run | Manual/UI/API | Chưa rerun explicit admin status filter verification as testcase in this report cycle. |
| REG-ORD-09 | Order | Pass | Automated | Video: [fe-admin--admin-can-review-order-timeline-and-update-status-with-reason--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-review-order-timeline-and-update-status-with-reason--passed.webm) |
| REG-ORD-10 | Order | Not Run | Manual/API | Chưa verify invalid status transition negative case. |
| REG-PAY-01 | Payment | Not Run | Manual/API | Chưa chạy explicit online payment creation regression. |
| REG-PAY-02 | Payment | Not Run | Manual/API | Chưa chạy live callback success regression. |
| REG-PAY-03 | Payment | Not Run | Manual/API | Chưa chạy live callback failed regression. |
| REG-PAY-04 | Payment | Not Run | Manual/API | Chưa verify idempotent duplicate callback. |
| REG-PAY-05 | Payment | Not Run | Manual/API/Job | Chưa verify expire pending payment job. |
| REG-TRK-01 | Tracking | Pass | Automated | Video: [fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm) |
| REG-TRK-02 | Tracking | Pass | Automated | Video: [fe-admin--admin-can-review-order-timeline-and-update-status-with-reason--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-review-order-timeline-and-update-status-with-reason--passed.webm) |
| REG-NOTI-01 | Notification | Pass | Automated | Video: [fe-customer--customer-can-view-notifications-and-mark-all-as-read--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-view-notifications-and-mark-all-as-read--passed.webm) |
| REG-NOTI-02 | Notification | Not Run | Manual/UI | Chưa verify mark single notification as read. |
| REG-NOTI-03 | Notification | Pass | Automated | Video: [fe-customer--customer-can-view-notifications-and-mark-all-as-read--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-view-notifications-and-mark-all-as-read--passed.webm) |
| REG-NOTI-04 | Notification | Not Run | Manual/API | Chưa verify notification dedupe on duplicate event. |
| REG-REV-01 | Reporting | Pass | Automated | Video: [fe-admin--admin-can-open-revenue-report--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-open-revenue-report--passed.webm) |
| REG-REV-02 | Reporting | Not Run | Manual/UI/API | Chưa verify multiple date-range combinations. |
| REG-REV-03 | Reporting | Not Run | Manual/API/DB | Chưa đối soát net revenue với source data. |
| REG-RBAC-01 | Security | Not Run | Manual/API | Chưa chạy explicit negative access user -> admin. |
| REG-RBAC-02 | Security | Not Run | Manual/API | Chưa chạy explicit admin misuse flow case. |
| REG-API-01 | Contract | Pass | API | Verified `POST /api/v1/auth/login` with `email/password`. |
| REG-API-02 | Contract | Not Run | API | Chưa kiểm full success/error envelope consistency matrix. |
| REG-OPS-01 | Smoke | Pass | API | Verified `GET /healthz`. |

## 4. Executed Evidence
- Playwright result:
  - `FE_admin`: `4 passed`
  - `FE_customer`: `4 passed`
- Video artifacts:
  - Admin login/dashboard: [fe-admin--admin-can-log-in-and-see-dashboard-data--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-log-in-and-see-dashboard-data--passed.webm)
  - Admin create product: [fe-admin--admin-can-create-a-product--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-create-a-product--passed.webm)
  - Admin order timeline/status: [fe-admin--admin-can-review-order-timeline-and-update-status-with-reason--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-review-order-timeline-and-update-status-with-reason--passed.webm)
  - Admin revenue report: [fe-admin--admin-can-open-revenue-report--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-open-revenue-report--passed.webm)
  - Customer login/browse: [fe-customer--customer-can-log-in-and-browse-products--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-log-in-and-browse-products--passed.webm)
  - Customer COD checkout: [fe-customer--customer-can-complete-cod-checkout--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-complete-cod-checkout--passed.webm)
  - Customer order detail/tracking: [fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm)
  - Customer notifications: [fe-customer--customer-can-view-notifications-and-mark-all-as-read--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-view-notifications-and-mark-all-as-read--passed.webm)
- API result:
  - `/healthz` returned `success=true`, `status=ok`
  - `/api/v1/auth/login` returned `200 OK` with token payload

## 5. Overall Assessment
- Regression subset đã chạy: `Pass`
- Không có lỗi phát hiện trong phần đã automated và smoke-tested.
- Report này chưa tương đương full regression sign-off vì còn nhiều case manual/API chưa chạy.

## 6. Remaining Work Before Full Regression Sign-off
- Chạy negative auth cases.
- Chạy payment callback/idempotency/expiry cases.
- Chạy product validation/edit/delete cases.
- Chạy order cancel / invalid transition negative cases.
- Chạy notification dedupe verification.
- Chạy revenue reconciliation và status filter regression theo dữ liệu thật.
- Chạy RBAC negative cases.
