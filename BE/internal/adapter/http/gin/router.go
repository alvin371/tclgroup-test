package gin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tclgroup/stock-management/internal/adapter/http/gin/handler"
	"github.com/tclgroup/stock-management/internal/adapter/http/gin/middleware"
	uc "github.com/tclgroup/stock-management/internal/usecase/inventory"
)

// SetupRouter wires all use cases into gin route handlers and returns the engine.
func SetupRouter(
	stockInUC *uc.StockInUseCase,
	stockOutUC *uc.StockOutUseCase,
	inventoryUC *uc.InventoryUseCase,
	reportUC *uc.ReportUseCase,
	log *zap.Logger,
) *gin.Engine {
	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))
	r.Use(middleware.Recovery(log))

	// Handlers
	stockInH := handler.NewStockInHandler(stockInUC)
	stockOutH := handler.NewStockOutHandler(stockOutUC)
	inventoryH := handler.NewInventoryHandler(inventoryUC)
	reportH := handler.NewReportHandler(reportUC)

	v1 := r.Group("/api/v1")
	{
		// Products
		products := v1.Group("/products")
		{
			products.POST("", inventoryH.CreateProduct)
			products.GET("", inventoryH.ListProducts)
		}

		// Stock In
		stockIn := v1.Group("/stock-in")
		{
			stockIn.POST("", stockInH.Create)
			stockIn.GET("", stockInH.List)
			stockIn.GET("/:id", stockInH.GetDetail)
			stockIn.PATCH("/:id/advance", stockInH.Advance)
			stockIn.DELETE("/:id", stockInH.Cancel)
		}

		// Stock Out
		stockOut := v1.Group("/stock-out")
		{
			stockOut.POST("/allocate", stockOutH.Allocate)
			stockOut.GET("", stockOutH.List)
			stockOut.GET("/:id", stockOutH.GetDetail)
			stockOut.PATCH("/:id/execute", stockOutH.Execute)
			stockOut.DELETE("/:id", stockOutH.Cancel)
		}

		// Inventory
		inventory := v1.Group("/inventory")
		{
			inventory.GET("", inventoryH.List)
			inventory.PATCH("/:product_id/adjust", inventoryH.Adjust)
		}

		// Reports
		reports := v1.Group("/reports")
		{
			reports.GET("/stock-in", reportH.StockInReport)
			reports.GET("/stock-out", reportH.StockOutReport)
		}
	}

	return r
}
