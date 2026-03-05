package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger returns a gin middleware that logs each request with zap.
func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		requestID := GetRequestID(c.Request.Context())

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Error("request error", append(fields, zap.Error(e.Err))...)
			}
		} else {
			log.Info("request", fields...)
		}
	}
}
