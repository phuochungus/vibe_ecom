package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"golf-store/be-mono/internal/shared/middleware"
)

type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func OK(c *gin.Context, data any) {
	JSON(c, http.StatusOK, data, nil)
}

func Created(c *gin.Context, data any) {
	JSON(c, http.StatusCreated, data, nil)
}

func Accepted(c *gin.Context, data any) {
	JSON(c, http.StatusAccepted, data, nil)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func JSON(c *gin.Context, status int, data any, meta any) {
	c.JSON(status, gin.H{
		"success":    true,
		"data":       data,
		"meta":       meta,
		"request_id": middleware.RequestIDFromContext(c),
		"timestamp":  time.Now().UTC().Format(time.RFC3339Nano),
	})
}

func Error(c *gin.Context, status int, code string, message string, details any) {
	c.JSON(status, gin.H{
		"success": false,
		"error": ErrorObject{
			Code:    code,
			Message: message,
			Details: details,
		},
		"request_id": middleware.RequestIDFromContext(c),
		"timestamp":  time.Now().UTC().Format(time.RFC3339Nano),
	})
}
