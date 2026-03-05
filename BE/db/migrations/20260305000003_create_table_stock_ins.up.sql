CREATE TABLE IF NOT EXISTS stock_ins (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id),
    quantity   BIGINT NOT NULL CHECK (quantity > 0),
    status     VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    notes      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
