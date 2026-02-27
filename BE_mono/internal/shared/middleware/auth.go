package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"golf-store/be-mono/internal/shared/model"
)

type TokenResolver interface {
	ResolveAccessToken(token string) (*model.User, bool)
}

const contextUserKey = "auth_user"

func RequireAuth(resolver TokenResolver) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := parseBearer(c.GetHeader("Authorization"))
		if token == "" {
			abortUnauthorized(c)
			c.Abort()
			return
		}

		user, ok := resolver.ResolveAccessToken(token)
		if !ok {
			abortUnauthorized(c)
			c.Abort()
			return
		}

		c.Set(contextUserKey, user)
		c.Next()
	}
}

func RequireRole(role model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := UserFromContext(c)
		if user == nil {
			abortUnauthorized(c)
			c.Abort()
			return
		}
		if user.Role != role {
			abortForbidden(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

func UserFromContext(c *gin.Context) *model.User {
	v, ok := c.Get(contextUserKey)
	if !ok {
		return nil
	}

	user, ok := v.(*model.User)
	if !ok {
		return nil
	}

	return user
}

func parseBearer(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": "Unauthorized",
		},
		"request_id": RequestIDFromContext(c),
		"timestamp":  time.Now().UTC().Format(time.RFC3339Nano),
	})
}

func abortForbidden(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "FORBIDDEN",
			"message": "Forbidden",
		},
		"request_id": RequestIDFromContext(c),
		"timestamp":  time.Now().UTC().Format(time.RFC3339Nano),
	})
}
