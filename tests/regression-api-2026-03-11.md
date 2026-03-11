# Live API Regression Evidence - 2026-03-11

Environment:
- Backend: `http://127.0.0.1:8080`
- Database: `golf-store-mono-postgres`
- Accounts:
  - `admin@golf.local`
  - `user@golf.local`

Results:
- `REG-ACC-03`: pass. Wrong password returned `401 INVALID_CREDENTIALS`.
- `REG-ACC-04`: pass. After five failed logins, the sixth attempt returned `423 ACCOUNT_LOCKED`.
- `REG-ACC-05`: pass. Refresh token returned a new access token.
- `REG-OPS-01`: pass. `/readyz` returned `200` with `ready=true`.
- `REG-API-02`: pass. Success and error envelopes both matched the standard response shape.
- `REG-PROD-03`: pass. Admin updated product price to `175000` and stock to `42`.
- `REG-PROD-04`: pass. Deleted product returned `204`, and the customer product endpoint returned `404`.
- `REG-PROD-05`: pass. Creating a product with `price=0` returned `400 VALIDATION_ERROR`.
- `REG-PROD-06`: pass. Creating a product with `stock=-1` returned `400 VALIDATION_ERROR`.
- `REG-PROD-07`: pass. The discontinued/deleted product was hidden from the customer product API.
- `REG-ORD-03`: pass. Two created orders produced different order codes.
- `REG-ORD-06`: pass. User canceled an eligible order and product stock was restored.
- `REG-ORD-07`: pass. Cancel after `SHIPPING` returned `409 ORDER_CANNOT_CANCEL`.
- `REG-ORD-08`: pass. Admin status filters for `PAID` and `SHIPPING` returned different order sets.
- `REG-ORD-10`: pass. Invalid status back-transition returned `409 CONFLICT`.
- `REG-PAY-01`: pass. Online payment creation returned a `checkout_url`.
- `REG-PAY-02`: pass. Success webhook updated the order to `order_status=PAID`, `payment_status=PAID`.
- `REG-PAY-03`: pass. Failed webhook updated `payment_status=FAILED` and kept `order_status=PENDING_PAYMENT`.
- `REG-PAY-04`: pass. Duplicate webhook was deduplicated and did not create a second transaction row.
- `REG-PAY-05`: pass. Expired pending-payment order auto-canceled and restored stock.
- `REG-NOTI-02`: pass. Mark-one-as-read returned `read=true`.
- `REG-NOTI-04`: pass. Duplicate payment success event produced only one notification row.
- `REG-RBAC-01`: pass. User token was blocked from admin API with `403 FORBIDDEN`.
- `REG-RBAC-02`: pass. Admin token calling customer notifications endpoint returned a controlled `200` response.
- `REG-REV-02`: pass. Revenue summary changed across date ranges as expected.
- `REG-REV-03`: pass. Revenue summary reconciled with DB:
  - `gross_revenue=150000`
  - `refund_amount=0`
  - `net_revenue=150000`

Supporting observations:
- `REG-ORD-08` specifically verified the issue the user reported earlier:
  - `status=PAID` returned a set containing the paid test order only.
  - `status=SHIPPING` returned a different set containing the shipping test order only.
- Notification dedupe was verified against PostgreSQL with:
  - `event_type='payment.succeeded'`
  - `event_key=<order_id>`
  - row count remained `1`
