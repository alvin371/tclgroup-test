package main

import (
	"github.com/tclgroup/stock-management/internal/adapter/http/gin"
	"github.com/tclgroup/stock-management/internal/adapter/postgres"
	pgRepo "github.com/tclgroup/stock-management/internal/adapter/postgres/repo"
	"github.com/tclgroup/stock-management/internal/pkg/clock"
	"github.com/tclgroup/stock-management/internal/pkg/config"
	"github.com/tclgroup/stock-management/internal/pkg/idgen"
	"github.com/tclgroup/stock-management/internal/pkg/logger"
	uc "github.com/tclgroup/stock-management/internal/usecase/inventory"
)

func main() {
	// Configuration
	cfg := config.MustLoad()

	// Logger
	log := logger.New(cfg.Logger)
	defer func() { _ = log.Sync() }()

	// Database
	pool := postgres.MustConnect(cfg.Database)
	defer pool.Close()

	// Infrastructure
	clk := clock.RealClock{}
	gen := idgen.UUIDGenerator{}
	txManager := postgres.NewTxManager(pool)

	// Repositories
	productRepo := pgRepo.NewProductRepo(pool)
	inventoryRepo := pgRepo.NewInventoryRepo(pool)
	stockInRepo := pgRepo.NewStockInRepo(pool)
	stockOutRepo := pgRepo.NewStockOutRepo(pool)
	reservationRepo := pgRepo.NewReservationRepo(pool)
	historyRepo := pgRepo.NewHistoryRepo(pool)

	// Use Cases
	stockInUC := uc.NewStockInUseCase(stockInRepo, inventoryRepo, historyRepo, txManager, clk, gen, productRepo)
	stockOutUC := uc.NewStockOutUseCase(stockOutRepo, inventoryRepo, reservationRepo, historyRepo, txManager, clk, gen, productRepo)
	inventoryUC := uc.NewInventoryUseCase(inventoryRepo, productRepo, historyRepo, txManager, clk, gen)
	reportUC := uc.NewReportUseCase(stockInRepo, stockOutRepo)

	// HTTP Server
	r := gin.SetupRouter(stockInUC, stockOutUC, inventoryUC, reportUC, log)

	log.Sugar().Infof("starting server on %s", cfg.App.Address)
	if err := r.Run(cfg.App.Address); err != nil {
		log.Sugar().Fatalf("server error: %v", err)
	}
}
