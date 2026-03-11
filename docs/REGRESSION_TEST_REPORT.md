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
  - API evidence: [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md)

## 2. Execution Summary
- Automated Playwright:
  - `FE_admin`: `4/4 passed`
  - `FE_customer`: `4/4 passed`
  - Video recording: `enabled`
- Live API/DB regression:
  - Additional covered cases: `26/26 passed`
- Smoke:
  - `GET /healthz`: `passed`
  - `GET /readyz`: `passed`
  - `POST /api/v1/auth/login`: `passed`

## 3. Test Case Result Matrix

| TC ID | Module | Status | Execution Type | Evidence / Notes |
|---|---|---|---|---|
| REG-ACC-01 | Auth | Pass | Automated UI | Video: [fe-customer--customer-can-log-in-and-browse-products--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-log-in-and-browse-products--passed.webm) |
| REG-ACC-02 | Auth | Pass | Automated UI | Video: [fe-admin--admin-can-log-in-and-see-dashboard-data--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-log-in-and-see-dashboard-data--passed.webm) |
| REG-ACC-03 | Auth | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-ACC-04 | Auth | Pass | API/DB | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-ACC-05 | Auth | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-PROD-01 | Product | Pass | Automated UI | Video: [fe-admin--admin-can-log-in-and-see-dashboard-data--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-log-in-and-see-dashboard-data--passed.webm) |
| REG-PROD-02 | Product | Pass | Automated UI | Video: [fe-admin--admin-can-create-a-product--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-create-a-product--passed.webm) |
| REG-PROD-03 | Product | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-PROD-04 | Product | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-PROD-05 | Product | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-PROD-06 | Product | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-PROD-07 | Product | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-ORD-01 | Order | Pass | Automated UI | Video: [fe-customer--customer-can-log-in-and-browse-products--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-log-in-and-browse-products--passed.webm) |
| REG-ORD-02 | Order | Pass | Automated UI | Video: [fe-customer--customer-can-complete-cod-checkout--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-complete-cod-checkout--passed.webm) |
| REG-ORD-03 | Order | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-ORD-04 | Order | Pass | Automated UI | Video: [fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm) |
| REG-ORD-05 | Order | Pass | Automated UI | Video: [fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm) |
| REG-ORD-06 | Order | Pass | API/DB | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-ORD-07 | Order | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-ORD-08 | Order | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-ORD-09 | Order | Pass | Automated UI | Video: [fe-admin--admin-can-review-order-timeline-and-update-status-with-reason--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-review-order-timeline-and-update-status-with-reason--passed.webm) |
| REG-ORD-10 | Order | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-PAY-01 | Payment | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-PAY-02 | Payment | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-PAY-03 | Payment | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-PAY-04 | Payment | Pass | API/DB | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-PAY-05 | Payment | Pass | API/DB | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-TRK-01 | Tracking | Pass | Automated UI | Video: [fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-view-order-detail-and-tracking-timeline--passed.webm) |
| REG-TRK-02 | Tracking | Pass | Automated UI | Video: [fe-admin--admin-can-review-order-timeline-and-update-status-with-reason--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-review-order-timeline-and-update-status-with-reason--passed.webm) |
| REG-NOTI-01 | Notification | Pass | Automated UI | Video: [fe-customer--customer-can-view-notifications-and-mark-all-as-read--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-view-notifications-and-mark-all-as-read--passed.webm) |
| REG-NOTI-02 | Notification | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-NOTI-03 | Notification | Pass | Automated UI | Video: [fe-customer--customer-can-view-notifications-and-mark-all-as-read--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-customer--customer-can-view-notifications-and-mark-all-as-read--passed.webm) |
| REG-NOTI-04 | Notification | Pass | API/DB | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-REV-01 | Reporting | Pass | Automated UI | Video: [fe-admin--admin-can-open-revenue-report--passed.webm](/Users/imacvip/vibe_ecom/tests/fe-admin--admin-can-open-revenue-report--passed.webm) |
| REG-REV-02 | Reporting | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-REV-03 | Reporting | Pass | API/DB | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-RBAC-01 | Security | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-RBAC-02 | Security | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-API-01 | Contract | Pass | API | Login API verified with `email/password`. |
| REG-API-02 | Contract | Pass | API | [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md) |
| REG-OPS-01 | Smoke | Pass | API | `/healthz` and `/readyz` both returned healthy responses. |

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
- API/DB evidence:
  - Live regression evidence: [regression-api-2026-03-11.md](/Users/imacvip/vibe_ecom/tests/regression-api-2026-03-11.md)
  - Includes auth lock, refresh, product validation, order cancel/filter, payment webhook/idempotency, notification dedupe, RBAC, revenue reconciliation, and `/readyz`

## 5. Overall Assessment
- Full regression suite for the current Sprint-01 scope was executed and passed.
- The specific admin order status filter concern was rechecked live and passed: `PAID` and `SHIPPING` returned different result sets.
- Current repo state is suitable for full regression sign-off for the documented Sprint-01 cases.

## 6. Remaining Work Before Full Regression Sign-off
- None for the documented `REG-*` suite in [REGRESSION_TEST_CASES.md](/Users/imacvip/vibe_ecom/docs/REGRESSION_TEST_CASES.md).
