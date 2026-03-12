-- 000001_init.up.sql
-- PostgreSQL 16
-- Enum types
CREATE TYPE user_role AS ENUM ('ADMIN', 'USER');

-- users
CREATE TABLE
    users (
        id VARCHAR(36) PRIMARY KEY,
        email VARCHAR(255) NOT NULL,
        phone VARCHAR(20) NOT NULL,
        password VARCHAR(255) NOT NULL,
        full_name VARCHAR(150) NOT NULL,
        role user_role NOT NULL DEFAULT 'USER',
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT uq_users_email UNIQUE (email),
        CONSTRAINT uq_users_phone UNIQUE (phone)
    );

-- products
CREATE TABLE
    products (
        id VARCHAR(36) PRIMARY KEY,
        sku VARCHAR(64) NOT NULL,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        price BIGINT NOT NULL,
        stock INTEGER NOT NULL,
        status VARCHAR(16) NOT NULL DEFAULT 'ACTIVE',
        image_url VARCHAR(500),
        image_urls JSONB,
        deleted_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT uq_products_sku UNIQUE (sku)
    );

-- orders
CREATE TABLE
    orders (
        id VARCHAR(36) PRIMARY KEY,
        order_code VARCHAR(32) NOT NULL,
        user_id VARCHAR(36) NOT NULL,
        order_status VARCHAR(32) NOT NULL,
        subtotal_amount BIGINT NOT NULL,
        discount_amount BIGINT NOT NULL,
        shipping_fee BIGINT NOT NULL,
        total_amount BIGINT NOT NULL,
        payment_due_at TIMESTAMPTZ,
        customer_note VARCHAR(500),
        cancel_reason VARCHAR(255),
        full_address VARCHAR(500) NOT NULL,
        placed_at TIMESTAMPTZ NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT uq_orders_order_code UNIQUE (order_code),
        CONSTRAINT uq_orders_user_idempotency UNIQUE (user_id, order_code),
        CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users (id)
    );

-- order_items
CREATE TABLE
    order_items (
        id VARCHAR(36) PRIMARY KEY,
        order_id VARCHAR(36) NOT NULL,
        product_id VARCHAR(36) NOT NULL,
        sku VARCHAR(64) NOT NULL,
        name VARCHAR(255) NOT NULL,
        unit_price BIGINT NOT NULL,
        quantity INTEGER NOT NULL,
        line_total BIGINT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT uq_order_item UNIQUE (order_id, product_id),
        CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders (id),
        CONSTRAINT fk_order_items_product FOREIGN KEY (product_id) REFERENCES products (id)
    );

-- user_addresses
CREATE TABLE
    user_addresses (
        id VARCHAR(36) PRIMARY KEY,
        user_id VARCHAR(36) NOT NULL,
        recipient_name VARCHAR(150) NOT NULL,
        recipient_phone VARCHAR(20) NOT NULL,
        full_address VARCHAR(500) NOT NULL,
        is_default BOOLEAN NOT NULL DEFAULT FALSE,
        deleted_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT fk_user_addresses_user FOREIGN KEY (user_id) REFERENCES users (id)
    );

-- notifications
CREATE TABLE
    notifications (
        id VARCHAR(36) PRIMARY KEY,
        user_id VARCHAR(36) NOT NULL,
        title VARCHAR(255) NOT NULL,
        content TEXT NOT NULL,
        status VARCHAR(16) NOT NULL DEFAULT 'SENT',
        is_read BOOLEAN NOT NULL DEFAULT FALSE,
        sent_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT uq_notification_dedupe UNIQUE (user_id, id),
        CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users (id)
    );

-- audit_logs
CREATE TABLE
    audit_logs (
        id VARCHAR(36) PRIMARY KEY,
        entity_type VARCHAR(64) NOT NULL,
        entity_id VARCHAR(36) NOT NULL,
        action VARCHAR(64) NOT NULL,
        actor_id VARCHAR(36),
        before_data JSONB,
        after_data JSONB,
        metadata JSONB,
        ip_address VARCHAR(45),
        user_agent VARCHAR(255),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- Indexes
-- user_addresses
CREATE INDEX idx_user_addresses__user_id ON user_addresses (user_id);

CREATE INDEX idx_user_addresses__user_default ON user_addresses (user_id, is_default);

-- audit_logs
CREATE INDEX idx_audit__entity_time ON audit_logs (entity_type, entity_id, created_at);

CREATE INDEX idx_audit__action_time ON audit_logs (action, created_at);