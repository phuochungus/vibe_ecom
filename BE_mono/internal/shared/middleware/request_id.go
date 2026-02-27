package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderRequestID = "X-Request-Id"
const requestIDContextKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(HeaderRequestID)
		if id == "" {
			id = uuid.NewString()
		}

		c.Set(requestIDContextKey, id)
		c.Writer.Header().Set(HeaderRequestID, id)
		c.Next()
	}
}

func RequestIDFromContext(c *gin.Context) string {
	v, ok := c.Get(requestIDContextKey)
	if !ok {
		return ""
	}

	id, ok := v.(string)
	if !ok {
		return ""
	}

	return id
}
