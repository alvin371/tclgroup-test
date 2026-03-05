package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/tclgroup/stock-management/internal/adapter/http/gin/dto"
	"github.com/tclgroup/stock-management/internal/adapter/http/gin/presenter"
	"github.com/tclgroup/stock-management/internal/pkg/httpx"
	"github.com/tclgroup/stock-management/internal/pkg/pagination"
	uc "github.com/tclgroup/stock-management/internal/usecase/inventory"
	ucDto "github.com/tclgroup/stock-management/internal/usecase/inventory/dto"
	"github.com/tclgroup/stock-management/internal/usecase/inventory/port"
)

// InventoryHandler handles HTTP requests for inventory and product operations.
type InventoryHandler struct {
	uc *uc.InventoryUseCase
}

// NewInventoryHandler creates a new InventoryHandler.
func NewInventoryHandler(uc *uc.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{uc: uc}
}

// CreateProduct handles POST /api/v1/products.
func (h *InventoryHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	out, err := h.uc.CreateProduct(c.Request.Context(), ucDto.CreateProductInput{
		SKU:        req.SKU,
		Name:       req.Name,
		CustomerID: req.CustomerID,
	})
	if err != nil {
		presenter.HandleError(c, err)
		return
	}

	httpx.Created(c, "Success create data.", out)
}

// ListProducts handles GET /api/v1/products.
func (h *InventoryHandler) ListProducts(c *gin.Context) {
	var pq dto.PaginationQuery
	if err := c.ShouldBindQuery(&pq); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	pg := pagination.New(pq.Page, pq.PerPage)
	out, err := h.uc.ListProducts(c.Request.Context(), port.ProductFilter{}, port.Pagination{Page: pg.Page, Limit: pg.Limit})
	if err != nil {
		presenter.HandleError(c, err)
		return
	}

	httpx.Paginated(c, out.Items, out.Total, pg.Page, pg.Limit, "Success get all data.")
}

// List handles GET /api/v1/inventory.
func (h *InventoryHandler) List(c *gin.Context) {
	var pq dto.PaginationQuery
	if err := c.ShouldBindQuery(&pq); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	pg := pagination.New(pq.Page, pq.PerPage)
	out, err := h.uc.List(c.Request.Context(), port.InventoryFilter{}, port.Pagination{Page: pg.Page, Limit: pg.Limit})
	if err != nil {
		presenter.HandleError(c, err)
		return
	}

	httpx.Paginated(c, out.Items, out.Total, pg.Page, pg.Limit, "Success get all data.")
}

// Adjust handles PATCH /api/v1/inventory/:product_id/adjust.
func (h *InventoryHandler) Adjust(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid product UUID")
		return
	}

	var req dto.AdjustStockRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err = h.uc.Adjust(c.Request.Context(), ucDto.AdjustStockInput{
		ProductID: productID,
		NewQty:    req.NewQty,
		Notes:     req.Notes,
	}); err != nil {
		presenter.HandleError(c, err)
		return
	}

	httpx.Success(c, "Success update data.", nil)
}
