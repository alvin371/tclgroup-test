CREATE INDEX IF NOT EXISTS idx_inventories_product_id ON inventories(product_id);
CREATE INDEX IF NOT EXISTS idx_stock_ins_product_id   ON stock_ins(product_id);
CREATE INDEX IF NOT EXISTS idx_stock_ins_status       ON stock_ins(status);
CREATE INDEX IF NOT EXISTS idx_stock_outs_product_id  ON stock_outs(product_id);
CREATE INDEX IF NOT EXISTS idx_stock_outs_status      ON stock_outs(status);
CREATE INDEX IF NOT EXISTS idx_reservations_product   ON reservations(product_id, status);
CREATE INDEX IF NOT EXISTS idx_reservations_stockout  ON reservations(stock_out_id);
CREATE INDEX IF NOT EXISTS idx_histories_entity       ON histories(entity_id, entity_type);
CREATE INDEX IF NOT EXISTS idx_products_sku           ON products(sku);
