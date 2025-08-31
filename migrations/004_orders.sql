-- +goose Up
-- +goose StatementBegin
-- Use pgcrypto's gen_random_uuid(); alternatively use "uuid-ossp" + uuid_generate_v4()
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'completed', 'canceled')),
    total_amount NUMERIC(10, 2) NOT NULL CHECK (total_amount >= 0),
    payment_intent_id VARCHAR(255) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_items;
-- +goose StatementEnd