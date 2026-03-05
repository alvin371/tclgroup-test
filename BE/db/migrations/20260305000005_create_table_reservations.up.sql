CREATE TABLE IF NOT EXISTS reservations (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stock_out_id UUID NOT NULL REFERENCES stock_outs(id),
    product_id   UUID NOT NULL REFERENCES products(id),
    quantity     BIGINT NOT NULL CHECK (quantity > 0),
    status       VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
