CREATE TABLE IF NOT EXISTS inventories (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id     UUID NOT NULL UNIQUE REFERENCES products(id),
    physical_stock BIGINT NOT NULL DEFAULT 0 CHECK (physical_stock >= 0),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
