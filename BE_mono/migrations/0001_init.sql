-- Core schema for BE_mono MySQL persistence.

CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(36) PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  phone VARCHAR(20) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  full_name VARCHAR(150) NOT NULL,
  role ENUM('USER','ADMIN') NOT NULL DEFAULT 'USER',
  status ENUM('ACTIVE','LOCKED','DISABLED') NOT NULL DEFAULT 'ACTIVE',
  failed_login_attempts INT NOT NULL DEFAULT 0,
  locked_until DATETIME(3) NULL,
  last_login_at DATETIME(3) NULL,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS auth_tokens (
  access_token VARCHAR(512) PRIMARY KEY,
  refresh_token VARCHAR(512) NOT NULL UNIQUE,
  user_id VARCHAR(36) NOT NULL,
  expires_at DATETIME(3) NOT NULL,
  created_at DATETIME(3) NOT NULL,
  CONSTRAINT fk_auth_tokens_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS products (
  id VARCHAR(36) PRIMARY KEY,
  sku VARCHAR(64) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  price_cents BIGINT NOT NULL,
  stock INT NOT NULL,
  status ENUM('ACTIVE','INACTIVE','DISCONTINUED') NOT NULL DEFAULT 'ACTIVE',
  image_url VARCHAR(500) NULL,
  deleted_at DATETIME(3) NULL,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS orders (
  id VARCHAR(36) PRIMARY KEY,
  order_code VARCHAR(32) NOT NULL UNIQUE,
  idempotency_key VARCHAR(128) NOT NULL,
  user_id VARCHAR(36) NOT NULL,
  order_status ENUM('NEW','PENDING_PAYMENT','PAID','PROCESSING','SHIPPING','COMPLETED','CANCELLED') NOT NULL,
  payment_status ENUM('UNPAID','PAID','FAILED','REFUNDED') NOT NULL,
  currency_code CHAR(3) NOT NULL,
  subtotal_amount BIGINT NOT NULL,
  discount_amount BIGINT NOT NULL,
  shipping_fee BIGINT NOT NULL,
  total_amount BIGINT NOT NULL,
  payment_due_at DATETIME(3) NULL,
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
  placed_at DATETIME(3) NOT NULL,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT uk_orders_user_idempotency UNIQUE (user_id, idempotency_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS order_items (
  id VARCHAR(36) PRIMARY KEY,
  order_id VARCHAR(36) NOT NULL,
  product_id VARCHAR(36) NOT NULL,
  sku VARCHAR(64) NOT NULL,
  name VARCHAR(255) NOT NULL,
  unit_price BIGINT NOT NULL,
  quantity INT NOT NULL,
  line_total BIGINT NOT NULL,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders(id),
  CONSTRAINT fk_order_items_product FOREIGN KEY (product_id) REFERENCES products(id),
  CONSTRAINT uk_order_item UNIQUE (order_id, product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS order_tracking_events (
  id VARCHAR(36) PRIMARY KEY,
  order_id VARCHAR(36) NOT NULL,
  from_status VARCHAR(32) NULL,
  to_status VARCHAR(32) NOT NULL,
  source_type VARCHAR(32) NOT NULL,
  description VARCHAR(500) NULL,
  occurred_at DATETIME(3) NOT NULL,
  CONSTRAINT fk_tracking_order FOREIGN KEY (order_id) REFERENCES orders(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS payment_transactions (
  id VARCHAR(36) PRIMARY KEY,
  order_id VARCHAR(36) NOT NULL,
  txn_type ENUM('PAYMENT','REFUND') NOT NULL DEFAULT 'PAYMENT',
  provider VARCHAR(50) NOT NULL,
  provider_txn_code VARCHAR(128) NULL UNIQUE,
  idempotency_key VARCHAR(128) NULL,
  amount BIGINT NOT NULL,
  currency_code CHAR(3) NOT NULL,
  status ENUM('PENDING','SUCCESS','FAILED','CANCELLED') NOT NULL,
  provider_response VARCHAR(500) NULL,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  CONSTRAINT fk_payment_order FOREIGN KEY (order_id) REFERENCES orders(id),
  CONSTRAINT uk_payment_order_idempotency UNIQUE (order_id, idempotency_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS notifications (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  channel ENUM('IN_APP','EMAIL','SMS') NOT NULL DEFAULT 'IN_APP',
  event_type VARCHAR(64) NOT NULL,
  event_key VARCHAR(128) NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  status ENUM('PENDING','SENT','FAILED') NOT NULL DEFAULT 'SENT',
  is_read TINYINT(1) NOT NULL DEFAULT 0,
  sent_at DATETIME(3) NULL,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT uk_notification_dedupe UNIQUE (user_id, event_type, event_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
