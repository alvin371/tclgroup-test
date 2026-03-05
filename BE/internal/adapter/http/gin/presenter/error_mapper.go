package presenter

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	domain "github.com/tclgroup/stock-management/internal/domain/inventory"
	"github.com/tclgroup/stock-management/internal/pkg/httpx"
)

type domainErrorEntry struct {
	status int
	code   string
}

var domainToHTTP = map[error]domainErrorEntry{
	domain.ErrProductNotFound:         {http.StatusNotFound, "PRODUCT_NOT_FOUND"},
	domain.ErrInventoryNotFound:        {http.StatusNotFound, "INVENTORY_NOT_FOUND"},
	domain.ErrStockInNotFound:          {http.StatusNotFound, "STOCK_IN_NOT_FOUND"},
	domain.ErrStockOutNotFound:         {http.StatusNotFound, "STOCK_OUT_NOT_FOUND"},
	domain.ErrReservationNotFound:      {http.StatusNotFound, "RESERVATION_NOT_FOUND"},
	domain.ErrInsufficientStock:        {http.StatusConflict, "INSUFFICIENT_STOCK"},
	domain.ErrInvalidStatusTransition:  {http.StatusUnprocessableEntity, "INVALID_STATUS_TRANSITION"},
	domain.ErrCannotCancelDoneStockIn:  {http.StatusConflict, "CANNOT_CANCEL_DONE_STOCK_IN"},
	domain.ErrInvalidQuantity:          {http.StatusBadRequest, "INVALID_QUANTITY"},
}

// HandleError maps a domain error to an HTTP response.
func HandleError(c *gin.Context, err error) {
	for domainErr, entry := range domainToHTTP {
		if errors.Is(err, domainErr) {
			httpx.Error(c, entry.status, entry.code, err.Error())
			return
		}
	}
	httpx.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
}
