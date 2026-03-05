package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type requestIDKey struct{}

const headerRequestID = "X-Request-ID"

// RequestID injects a request ID into the context and response headers.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(headerRequestID)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Header(headerRequestID, requestID)
		ctx := context.WithValue(c.Request.Context(), requestIDKey{}, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Set(headerRequestID, requestID)
		c.Next()
	}
}

// GetRequestID retrieves the request ID from context.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}
