package observability

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderCorrelationID = "X-Correlation-ID"
const correlationContextKey = "correlation_id"

func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(HeaderCorrelationID)
		if id == "" {
			id = uuid.NewString()
		}

		c.Set(correlationContextKey, id)
		c.Writer.Header().Set(HeaderCorrelationID, id)
		c.Next()
	}
}

func CorrelationIDFromContext(c *gin.Context) string {
	v, ok := c.Get(correlationContextKey)
	if !ok {
		return ""
	}

	id, ok := v.(string)
	if !ok {
		return ""
	}

	return id
}
