-- Core schema for BE_mono PostgreSQL persistence.

CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(36) PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  phone VARCHAR(20) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  full_name VARCHAR(150) NOT NULL,
  role VARCHAR(16) NOT NULL DEFAULT 'USER',
  status VARCHAR(16) NOT NULL DEFAULT 'ACTIVE',
  failed_login_attempts INTEGER NOT NULL DEFAULT 0,
  locked_until TIMESTAMPTZ NULL,
  last_login_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS auth_tokens (
  access_token VARCHAR(512) PRIMARY KEY,
  refresh_token VARCHAR(512) NOT NULL UNIQUE,
  user_id VARCHAR(36) NOT NULL REFERENCES users(id),
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
  id VARCHAR(36) PRIMARY KEY,
  sku VARCHAR(64) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  price BIGINT NOT NULL,
  stock INTEGER NOT NULL,
  status VARCHAR(16) NOT NULL DEFAULT 'ACTIVE',
  image_url VARCHAR(500) NULL,
  image_urls JSONB NULL,
  deleted_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
  id VARCHAR(36) PRIMARY KEY,
  order_code VARCHAR(32) NOT NULL UNIQUE,
  idempotency_key VARCHAR(128) NOT NULL,
  user_id VARCHAR(36) NOT NULL REFERENCES users(id),
  order_status VARCHAR(32) NOT NULL,
  payment_status VARCHAR(16) NOT NULL,
  currency_code CHAR(3) NOT NULL,
  subtotal_amount BIGINT NOT NULL,
  discount_amount BIGINT NOT NULL,
  shipping_fee BIGINT NOT NULL,
  total_amount BIGINT NOT NULL,
  payment_due_at TIMESTAMPTZ NULL,
  customer_note VARCHAR(500) NULL,
  cancel_reason VARCHAR(255) NULL,
  shipping_recipient_name VARCHAR(150) NOT NULL,
  shipping_phone VARCHAR(20) NOT NULL,
  shipping_line1 VARCHAR(255) NOT NULL,
  shipping_line2 VARCHAR(255) NULL,
  shipping_ward VARCHAR(120) NULL,
  shipping_district VARCHAR(120) NULL,
  shipping_city VARCHAR(120) NOT NULL,
  shipping_province VARCHAR(120) NULL,
  shipping_postal_code VARCHAR(20) NULL,
  shipping_country_code CHAR(2) NOT NULL,
  placed_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT uk_orders_user_idempotency UNIQUE (user_id, idempotency_key)
);

CREATE TABLE IF NOT EXISTS order_items (
  id VARCHAR(36) PRIMARY KEY,
  order_id VARCHAR(36) NOT NULL REFERENCES orders(id),
  product_id VARCHAR(36) NOT NULL REFERENCES products(id),
  sku VARCHAR(64) NOT NULL,
  name VARCHAR(255) NOT NULL,
  unit_price BIGINT NOT NULL,
  quantity INTEGER NOT NULL,
  line_total BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT uk_order_item UNIQUE (order_id, product_id)
);

CREATE TABLE IF NOT EXISTS order_tracking_events (
  id VARCHAR(36) PRIMARY KEY,
  order_id VARCHAR(36) NOT NULL REFERENCES orders(id),
  from_status VARCHAR(32) NULL,
  to_status VARCHAR(32) NOT NULL,
  source_type VARCHAR(32) NOT NULL,
  description VARCHAR(500) NULL,
  occurred_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS payment_transactions (
  id VARCHAR(36) PRIMARY KEY,
  order_id VARCHAR(36) NOT NULL REFERENCES orders(id),
  txn_type VARCHAR(16) NOT NULL DEFAULT 'PAYMENT',
  provider VARCHAR(50) NOT NULL,
  provider_txn_code VARCHAR(128) NULL UNIQUE,
  idempotency_key VARCHAR(128) NULL,
  amount BIGINT NOT NULL,
  currency_code CHAR(3) NOT NULL,
  status VARCHAR(16) NOT NULL,
  provider_response VARCHAR(500) NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT uk_payment_order_idempotency UNIQUE (order_id, idempotency_key)
);

CREATE TABLE IF NOT EXISTS notifications (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL REFERENCES users(id),
  channel VARCHAR(16) NOT NULL DEFAULT 'IN_APP',
  event_type VARCHAR(64) NOT NULL,
  event_key VARCHAR(128) NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  status VARCHAR(16) NOT NULL DEFAULT 'SENT',
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  sent_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT uk_notification_dedupe UNIQUE (user_id, event_type, event_key)
);
