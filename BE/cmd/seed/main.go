// Package main provides a database seeder for development and testing.
// Run: go run ./cmd/seed
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclgroup/stock-management/internal/pkg/config"
)

func main() {
	cfg := config.MustLoad()
	dsn := cfg.Database

	pool, err := pgxpool.New(context.Background(), dsn.DSN)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	if err = pool.Ping(context.Background()); err != nil {
		log.Fatalf("ping: %v", err)
	}

	ctx := context.Background()

	steps := []struct {
		name string
		fn   func(context.Context, *pgxpool.Pool) error
	}{
		{"products", seedProducts},
		{"inventories", seedInventories},
		{"stock_ins", seedStockIns},
		{"stock_outs", seedStockOuts},
		{"reservations", seedReservations},
		{"histories", seedHistories},
	}

	for _, s := range steps {
		log.Printf("seeding %s ...", s.name)
		if err = s.fn(ctx, pool); err != nil {
			log.Fatalf("seed %s: %v", s.name, err)
		}
		log.Printf("seeding %s done", s.name)
	}

	log.Println("all seeds applied successfully")
}

// ─── Products ────────────────────────────────────────────────────────────────
// 10 products across 3 customers.
func seedProducts(ctx context.Context, pool *pgxpool.Pool) error {
	rows := []struct {
		id, sku, name, customerID string
	}{
		{"11000000-0000-0000-0000-000000000001", "SKU-LAPTOP-001", "Laptop Pro 15", "99000000-0000-0000-0000-000000000001"},
		{"11000000-0000-0000-0000-000000000002", "SKU-MOUSE-001", "Wireless Mouse", "99000000-0000-0000-0000-000000000001"},
		{"11000000-0000-0000-0000-000000000003", "SKU-KEYB-001", "Mechanical Keyboard", "99000000-0000-0000-0000-000000000001"},
		{"11000000-0000-0000-0000-000000000004", "SKU-MONIT-001", "4K Monitor 27\"", "99000000-0000-0000-0000-000000000001"},
		{"11000000-0000-0000-0000-000000000005", "SKU-HDPH-001", "Noise Cancelling Headphones", "99000000-0000-0000-0000-000000000002"},
		{"11000000-0000-0000-0000-000000000006", "SKU-WEBCAM-001", "HD Webcam 1080p", "99000000-0000-0000-0000-000000000002"},
		{"11000000-0000-0000-0000-000000000007", "SKU-CHAIR-001", "Ergonomic Office Chair", "99000000-0000-0000-0000-000000000002"},
		{"11000000-0000-0000-0000-000000000008", "SKU-DESK-001", "Standing Desk Adjustable", "99000000-0000-0000-0000-000000000003"},
		{"11000000-0000-0000-0000-000000000009", "SKU-CABLE-001", "USB-C Cable 2m", "99000000-0000-0000-0000-000000000003"},
		{"11000000-0000-0000-0000-000000000010", "SKU-HUB-001", "USB-C Hub 7-in-1", "99000000-0000-0000-0000-000000000003"},
	}

	for _, r := range rows {
		_, err := pool.Exec(ctx, `
			INSERT INTO products (id, sku, name, customer_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING`,
			r.id, r.sku, r.name, r.customerID,
		)
		if err != nil {
			return fmt.Errorf("insert product %s: %w", r.sku, err)
		}
	}
	return nil
}

// ─── Inventories ─────────────────────────────────────────────────────────────
// One inventory row per product with varying stock levels.
func seedInventories(ctx context.Context, pool *pgxpool.Pool) error {
	rows := []struct {
		id, productID string
		physicalStock int
	}{
		{"22000000-0000-0000-0000-000000000001", "11000000-0000-0000-0000-000000000001", 150},
		{"22000000-0000-0000-0000-000000000002", "11000000-0000-0000-0000-000000000002", 320},
		{"22000000-0000-0000-0000-000000000003", "11000000-0000-0000-0000-000000000003", 85},
		{"22000000-0000-0000-0000-000000000004", "11000000-0000-0000-0000-000000000004", 40},
		{"22000000-0000-0000-0000-000000000005", "11000000-0000-0000-0000-000000000005", 200},
		{"22000000-0000-0000-0000-000000000006", "11000000-0000-0000-0000-000000000006", 60},
		{"22000000-0000-0000-0000-000000000007", "11000000-0000-0000-0000-000000000007", 25},
		{"22000000-0000-0000-0000-000000000008", "11000000-0000-0000-0000-000000000008", 12},
		{"22000000-0000-0000-0000-000000000009", "11000000-0000-0000-0000-000000000009", 500},
		{"22000000-0000-0000-0000-000000000010", "11000000-0000-0000-0000-000000000010", 75},
	}

	for _, r := range rows {
		_, err := pool.Exec(ctx, `
			INSERT INTO inventories (id, product_id, physical_stock, created_at, updated_at)
			VALUES ($1, $2, $3, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING`,
			r.id, r.productID, r.physicalStock,
		)
		if err != nil {
			return fmt.Errorf("insert inventory for product %s: %w", r.productID, err)
		}
	}
	return nil
}

// ─── Stock Ins ────────────────────────────────────────────────────────────────
// 11 stock-in records across all statuses.
func seedStockIns(ctx context.Context, pool *pgxpool.Pool) error {
	rows := []struct {
		id, productID string
		qty           int
		status, notes string
	}{
		{"33000000-0000-0000-0000-000000000001", "11000000-0000-0000-0000-000000000001", 50, "DONE", "Initial laptop stock"},
		{"33000000-0000-0000-0000-000000000002", "11000000-0000-0000-0000-000000000001", 100, "DONE", "Restock laptops Q1"},
		{"33000000-0000-0000-0000-000000000003", "11000000-0000-0000-0000-000000000002", 200, "DONE", "Mouse bulk purchase"},
		{"33000000-0000-0000-0000-000000000004", "11000000-0000-0000-0000-000000000003", 50, "DONE", "Keyboard initial stock"},
		{"33000000-0000-0000-0000-000000000005", "11000000-0000-0000-0000-000000000004", 30, "IN_PROGRESS", "Monitor delivery in transit"},
		{"33000000-0000-0000-0000-000000000006", "11000000-0000-0000-0000-000000000005", 120, "DONE", "Headphone bulk order"},
		{"33000000-0000-0000-0000-000000000007", "11000000-0000-0000-0000-000000000006", 40, "CREATED", "Webcam pending receipt"},
		{"33000000-0000-0000-0000-000000000008", "11000000-0000-0000-0000-000000000007", 20, "CANCELLED", "Chair order cancelled"},
		{"33000000-0000-0000-0000-000000000009", "11000000-0000-0000-0000-000000000008", 10, "DONE", "Desk delivery complete"},
		{"33000000-0000-0000-0000-000000000010", "11000000-0000-0000-0000-000000000009", 300, "DONE", "Cable bulk stock"},
		{"33000000-0000-0000-0000-000000000011", "11000000-0000-0000-0000-000000000010", 60, "IN_PROGRESS", "Hub shipment en route"},
	}

	for _, r := range rows {
		_, err := pool.Exec(ctx, `
			INSERT INTO stock_ins (id, product_id, quantity, status, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING`,
			r.id, r.productID, r.qty, r.status, r.notes,
		)
		if err != nil {
			return fmt.Errorf("insert stock_in %s: %w", r.id, err)
		}
	}
	return nil
}

// ─── Stock Outs ───────────────────────────────────────────────────────────────
// 10 stock-out records across all statuses.
func seedStockOuts(ctx context.Context, pool *pgxpool.Pool) error {
	rows := []struct {
		id, productID string
		qty           int
		status, notes string
	}{
		{"44000000-0000-0000-0000-000000000001", "11000000-0000-0000-0000-000000000001", 5, "DONE", "Laptop sale order #1001"},
		{"44000000-0000-0000-0000-000000000002", "11000000-0000-0000-0000-000000000001", 10, "DONE", "Laptop sale order #1002"},
		{"44000000-0000-0000-0000-000000000003", "11000000-0000-0000-0000-000000000002", 30, "DONE", "Mouse sale order #1003"},
		{"44000000-0000-0000-0000-000000000004", "11000000-0000-0000-0000-000000000003", 15, "DONE", "Keyboard sale order #1004"},
		{"44000000-0000-0000-0000-000000000005", "11000000-0000-0000-0000-000000000004", 5, "IN_PROGRESS", "Monitor order #1005 processing"},
		{"44000000-0000-0000-0000-000000000006", "11000000-0000-0000-0000-000000000005", 20, "DRAFT", "Headphone pending allocation"},
		{"44000000-0000-0000-0000-000000000007", "11000000-0000-0000-0000-000000000006", 8, "CANCELLED", "Webcam order cancelled"},
		{"44000000-0000-0000-0000-000000000008", "11000000-0000-0000-0000-000000000007", 3, "DONE", "Chair sale order #1008"},
		{"44000000-0000-0000-0000-000000000009", "11000000-0000-0000-0000-000000000009", 50, "DRAFT", "Cable bulk order pending"},
		{"44000000-0000-0000-0000-000000000010", "11000000-0000-0000-0000-000000000010", 10, "DONE", "Hub sale order #1010"},
	}

	for _, r := range rows {
		_, err := pool.Exec(ctx, `
			INSERT INTO stock_outs (id, product_id, quantity, status, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING`,
			r.id, r.productID, r.qty, r.status, r.notes,
		)
		if err != nil {
			return fmt.Errorf("insert stock_out %s: %w", r.id, err)
		}
	}
	return nil
}

// ─── Reservations ─────────────────────────────────────────────────────────────
// One reservation per stock-out; status mirrors the stock-out lifecycle.
func seedReservations(ctx context.Context, pool *pgxpool.Pool) error {
	rows := []struct {
		id, stockOutID, productID string
		qty                       int
		status                    string
	}{
		// CONSUMED — mirrors DONE stock-outs
		{"55000000-0000-0000-0000-000000000001", "44000000-0000-0000-0000-000000000001", "11000000-0000-0000-0000-000000000001", 5, "CONSUMED"},
		{"55000000-0000-0000-0000-000000000002", "44000000-0000-0000-0000-000000000002", "11000000-0000-0000-0000-000000000001", 10, "CONSUMED"},
		{"55000000-0000-0000-0000-000000000003", "44000000-0000-0000-0000-000000000003", "11000000-0000-0000-0000-000000000002", 30, "CONSUMED"},
		{"55000000-0000-0000-0000-000000000004", "44000000-0000-0000-0000-000000000004", "11000000-0000-0000-0000-000000000003", 15, "CONSUMED"},
		// ACTIVE — mirrors IN_PROGRESS / DRAFT stock-outs
		{"55000000-0000-0000-0000-000000000005", "44000000-0000-0000-0000-000000000005", "11000000-0000-0000-0000-000000000004", 5, "ACTIVE"},
		{"55000000-0000-0000-0000-000000000006", "44000000-0000-0000-0000-000000000006", "11000000-0000-0000-0000-000000000005", 20, "ACTIVE"},
		// RELEASED — mirrors CANCELLED stock-out
		{"55000000-0000-0000-0000-000000000007", "44000000-0000-0000-0000-000000000007", "11000000-0000-0000-0000-000000000006", 8, "RELEASED"},
		// CONSUMED — more DONE
		{"55000000-0000-0000-0000-000000000008", "44000000-0000-0000-0000-000000000008", "11000000-0000-0000-0000-000000000007", 3, "CONSUMED"},
		// ACTIVE — DRAFT
		{"55000000-0000-0000-0000-000000000009", "44000000-0000-0000-0000-000000000009", "11000000-0000-0000-0000-000000000009", 50, "ACTIVE"},
		// CONSUMED — DONE
		{"55000000-0000-0000-0000-000000000010", "44000000-0000-0000-0000-000000000010", "11000000-0000-0000-0000-000000000010", 10, "CONSUMED"},
	}

	for _, r := range rows {
		_, err := pool.Exec(ctx, `
			INSERT INTO reservations (id, stock_out_id, product_id, quantity, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING`,
			r.id, r.stockOutID, r.productID, r.qty, r.status,
		)
		if err != nil {
			return fmt.Errorf("insert reservation %s: %w", r.id, err)
		}
	}
	return nil
}

// ─── Histories ────────────────────────────────────────────────────────────────
// Audit trail for stock-in advances and stock-out executions.
func seedHistories(ctx context.Context, pool *pgxpool.Pool) error {
	rows := []struct {
		id, entityType, entityID string
		action                   string
		oldStatus, newStatus     string
		qty                      int
	}{
		// stock_in lifecycle for order 001
		{"66000000-0000-0000-0000-000000000001", "stock_in", "33000000-0000-0000-0000-000000000001", "advance", "CREATED", "IN_PROGRESS", 50},
		{"66000000-0000-0000-0000-000000000002", "stock_in", "33000000-0000-0000-0000-000000000001", "advance", "IN_PROGRESS", "DONE", 50},
		// stock_in lifecycle for order 002
		{"66000000-0000-0000-0000-000000000003", "stock_in", "33000000-0000-0000-0000-000000000002", "advance", "CREATED", "IN_PROGRESS", 100},
		{"66000000-0000-0000-0000-000000000004", "stock_in", "33000000-0000-0000-0000-000000000002", "advance", "IN_PROGRESS", "DONE", 100},
		// cancelled stock_in
		{"66000000-0000-0000-0000-000000000005", "stock_in", "33000000-0000-0000-0000-000000000008", "cancel", "CREATED", "CANCELLED", 20},
		// stock_out lifecycle for sale 001
		{"66000000-0000-0000-0000-000000000006", "stock_out", "44000000-0000-0000-0000-000000000001", "execute", "DRAFT", "IN_PROGRESS", 5},
		{"66000000-0000-0000-0000-000000000007", "stock_out", "44000000-0000-0000-0000-000000000001", "execute", "IN_PROGRESS", "DONE", 5},
		// stock_out lifecycle for sale 002
		{"66000000-0000-0000-0000-000000000008", "stock_out", "44000000-0000-0000-0000-000000000002", "execute", "DRAFT", "IN_PROGRESS", 10},
		{"66000000-0000-0000-0000-000000000009", "stock_out", "44000000-0000-0000-0000-000000000002", "execute", "IN_PROGRESS", "DONE", 10},
		// cancelled stock_out
		{"66000000-0000-0000-0000-000000000010", "stock_out", "44000000-0000-0000-0000-000000000007", "execute", "DRAFT", "CANCELLED", 8},
		// inventory adjustment
		{"66000000-0000-0000-0000-000000000011", "inventory", "22000000-0000-0000-0000-000000000001", "adjust", "150", "200", 200},
	}

	for _, r := range rows {
		_, err := pool.Exec(ctx, `
			INSERT INTO histories (id, entity_type, entity_id, action, old_status, new_status, quantity, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
			ON CONFLICT (id) DO NOTHING`,
			r.id, r.entityType, r.entityID, r.action, r.oldStatus, r.newStatus, r.qty,
		)
		if err != nil {
			return fmt.Errorf("insert history %s: %w", r.id, err)
		}
	}
	return nil
}
