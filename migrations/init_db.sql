-- migrate:up

CREATE TABLE IF NOT EXISTS orders (
    order_uid TEXT PRIMARY KEY,
    track_number TEXT,
    entry TEXT,
    delivery JSONB,
    payment JSONB,
    items JSONB,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INTEGER,
    date_created TIMESTAMP,
    oof_shard TEXT
);

-- migrate:down

DROP TABLE IF EXISTS orders;