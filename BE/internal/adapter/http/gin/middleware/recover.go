package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tclgroup/stock-management/internal/pkg/httpx"
)

// Recovery returns a gin middleware that recovers from panics and logs them.
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(c.Request.Context())
				log.Error("panic recovered",
					zap.Any("error", err),
					zap.String("request_id", requestID),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
				)
				httpx.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}
