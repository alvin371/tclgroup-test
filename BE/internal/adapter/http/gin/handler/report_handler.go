package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tclgroup/stock-management/internal/adapter/http/gin/dto"
	"github.com/tclgroup/stock-management/internal/adapter/http/gin/presenter"
	"github.com/tclgroup/stock-management/internal/pkg/httpx"
	"github.com/tclgroup/stock-management/internal/pkg/pagination"
	uc "github.com/tclgroup/stock-management/internal/usecase/inventory"
	"github.com/tclgroup/stock-management/internal/usecase/inventory/port"
)

// ReportHandler handles HTTP requests for report queries.
type ReportHandler struct {
	uc *uc.ReportUseCase
}

// NewReportHandler creates a new ReportHandler.
func NewReportHandler(uc *uc.ReportUseCase) *ReportHandler {
	return &ReportHandler{uc: uc}
}

// StockInReport handles GET /api/v1/reports/stock-in.
func (h *ReportHandler) StockInReport(c *gin.Context) {
	var pq dto.PaginationQuery
	if err := c.ShouldBindQuery(&pq); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	pg := pagination.New(pq.Page, pq.PerPage)
	out, err := h.uc.StockInReport(c.Request.Context(), port.StockInFilter{}, port.Pagination{Page: pg.Page, Limit: pg.Limit})
	if err != nil {
		presenter.HandleError(c, err)
		return
	}

	httpx.Paginated(c, out.Items, out.Total, pg.Page, pg.Limit, "Success get all data.")
}

// StockOutReport handles GET /api/v1/reports/stock-out.
func (h *ReportHandler) StockOutReport(c *gin.Context) {
	var pq dto.PaginationQuery
	if err := c.ShouldBindQuery(&pq); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	pg := pagination.New(pq.Page, pq.PerPage)
	out, err := h.uc.StockOutReport(c.Request.Context(), port.StockOutFilter{}, port.Pagination{Page: pg.Page, Limit: pg.Limit})
	if err != nil {
		presenter.HandleError(c, err)
		return
	}

	httpx.Paginated(c, out.Items, out.Total, pg.Page, pg.Limit, "Success get all data.")
}
