package observability

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCorrelationMiddlewareSetsHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(CorrelationID())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, CorrelationIDFromContext(c))
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	headerID := w.Header().Get(HeaderCorrelationID)
	if headerID == "" {
		t.Fatalf("expected correlation id header")
	}

	bodyID := strings.TrimSpace(w.Body.String())
	if bodyID == "" {
		t.Fatalf("expected correlation id in response body")
	}
	if bodyID != headerID {
		t.Fatalf("expected body correlation id to match header")
	}
}
