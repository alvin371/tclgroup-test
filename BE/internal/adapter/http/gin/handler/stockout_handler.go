package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/tclgroup/stock-management/internal/adapter/http/gin/dto"
	"github.com/tclgroup/stock-management/internal/adapter/http/gin/presenter"
	"github.com/tclgroup/stock-management/internal/domain/inventory/valueobject"
	"github.com/tclgroup/stock-management/internal/pkg/httpx"
	"github.com/tclgroup/stock-management/internal/pkg/pagination"
	uc "github.com/tclgroup/stock-management/internal/usecase/inventory"
	ucDto "github.com/tclgroup/stock-management/internal/usecase/inventory/dto"
	"github.com/tclgroup/stock-management/internal/usecase/inventory/port"
)

// StockOutHandler handles HTTP requests for stock-out operations.
type StockOutHandler struct {
	uc *uc.StockOutUseCase
}

// NewStockOutHandler creates a new StockOutHandler.
func NewStockOutHandler(uc *uc.StockOutUseCase) *StockOutHandler {
	return &StockOutHandler{uc: uc}
}

// Allocate handles POST /api/v1/stock-out/allocate.
func (h *StockOutHandler) Allocate(c *gin.Context) {
	var req dto.AllocateStockOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	out, err := h.uc.Allocate(c.Request.Context(), ucDto.AllocateStockOutInput{
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
		Notes:       req.Notes,
		UnitPrice:   req.UnitPrice,
		PerformedBy: req.PerformedBy,
		Location:    req.Location,
	})
	if err != nil {
		presenter.HandleError(c, err)
		return
	}

	httpx.Created(c, "Success create data.", out)
}

// Execute handles PATCH /api/v1/stock-out/:id/execute.
func (h *StockOutHandler) Execute(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid UUID")
		return
	}

	var req dto.ExecuteStockOutRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	out, err := h.uc.Execute(c.Request.Context(), id, valueobject.StockOutStatus(req.Status))
	if err != nil {
		presenter.HandleError(c, err)
		return
	}

	httpx.Success(c, "Success update data.", out)
}

// GetDetail handles GET /api/v1/stock-out/:id.
func (h *StockOutHandler) GetDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid UUID")
		return
	}
	out, err := h.uc.GetDetail(c.Request.Context(), id)
	if err != nil {
		presenter.HandleError(c, err)
		return
	}
	httpx.Success(c, "Success get data.", out)
}

// Cancel handles DELETE /api/v1/stock-out/:id.
func (h *StockOutHandler) Cancel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid UUID")
		return
	}

	if err = h.uc.Cancel(c.Request.Context(), id); err != nil {
		presenter.HandleError(c, err)
		return
	}

	httpx.Success(c, "Success delete data.", nil)
}

// List handles GET /api/v1/stock-out.
func (h *StockOutHandler) List(c *gin.Context) {
	var pq dto.PaginationQuery
	if err := c.ShouldBindQuery(&pq); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	pg := pagination.New(pq.Page, pq.PerPage)

	var filter port.StockOutFilter
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	out, err := h.uc.List(c.Request.Context(), filter, port.Pagination{Page: pg.Page, Limit: pg.Limit})
	if err != nil {
		presenter.HandleError(c, err)
		return
	}

	httpx.Paginated(c, out.Items, out.Total, pg.Page, pg.Limit, "Success get all data.")
}
