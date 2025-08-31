-- +goose Up
-- +goose StatementBegin
-- Use pgcrypto's gen_random_uuid(); alternatively use "uuid-ossp" + uuid_generate_v4()
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    book_id VARCHAR(64) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    price_cents INT NOT NULL CHECK (price_cents >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_items;
-- +goose StatementEnd