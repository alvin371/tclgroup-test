package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/tclgroup/stock-management/internal/adapter/http/gin/dto"
	"github.com/tclgroup/stock-management/internal/adapter/http/gin/presenter"
	"github.com/tclgroup/stock-management/internal/domain/inventory/valueobject"
	ucDto "github.com/tclgroup/stock-management/internal/usecase/inventory/dto"
	"github.com/tclgroup/stock-management/internal/usecase/inventory/port"
	"github.com/tclgroup/stock-management/internal/pkg/httpx"
	"github.com/tclgroup/stock-management/internal/pkg/pagination"
	uc "github.com/tclgroup/stock-management/internal/usecase/inventory"
)

// StockInHandler handles HTTP requests for stock-in operations.
type StockInHandler struct {
	uc *uc.StockInUseCase
}

// NewStockInHandler creates a new StockInHandler.
func NewStockInHandler(uc *uc.StockInUseCase) *StockInHandler {
	return &StockInHandler{uc: uc}
}

// Create handles POST /api/v1/stock-in.
func (h *StockInHandler) Create(c *gin.Context) {
	var req dto.CreateStockInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	out, err := h.uc.Create(c.Request.Context(), ucDto.CreateStockInInput{
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

// List handles GET /api/v1/stock-in.
func (h *StockInHandler) List(c *gin.Context) {
	var pq dto.PaginationQuery
	if err := c.ShouldBindQuery(&pq); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	pg := pagination.New(pq.Page, pq.PerPage)

	var filter port.StockInFilter
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

// Advance handles PATCH /api/v1/stock-in/:id/advance.
func (h *StockInHandler) Advance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid UUID")
		return
	}

	var req dto.AdvanceStockInRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	out, err := h.uc.Advance(c.Request.Context(), id, valueobject.StockInStatus(req.Status))
	if err != nil {
		presenter.HandleError(c, err)
		return
	}

	httpx.Success(c, "Success update data.", out)
}

// GetDetail handles GET /api/v1/stock-in/:id.
func (h *StockInHandler) GetDetail(c *gin.Context) {
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

// Cancel handles DELETE /api/v1/stock-in/:id.
func (h *StockInHandler) Cancel(c *gin.Context) {
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
