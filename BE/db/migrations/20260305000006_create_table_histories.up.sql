CREATE TABLE IF NOT EXISTS histories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(50)  NOT NULL,
    entity_id   UUID         NOT NULL,
    action      VARCHAR(100) NOT NULL,
    old_status  VARCHAR(50),
    new_status  VARCHAR(50),
    quantity    BIGINT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
