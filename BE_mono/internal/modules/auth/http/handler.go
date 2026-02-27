package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	authsvc "golf-store/be-mono/internal/modules/auth/service"
	"golf-store/be-mono/internal/shared/middleware"
	"golf-store/be-mono/internal/shared/model"
	"golf-store/be-mono/internal/shared/response"
)

type Handler struct {
	auth *authsvc.Service
}

func New(auth *authsvc.Service) *Handler {
	return &Handler{auth: auth}
}

type loginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) RegisterPublic(rg *gin.RouterGroup) {
	rg.POST("/auth/login", h.Login)
	rg.POST("/auth/refresh", h.RefreshToken)
}

func (h *Handler) RegisterProtected(rg *gin.RouterGroup) {
	rg.POST("/auth/logout", h.Logout)
	rg.GET("/me", h.Me)
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	result, apiErr := h.auth.Login(req.Identifier, req.Password)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}

	response.OK(c, gin.H{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"token_type":    result.TokenType,
		"expires_in":    result.ExpiresIn,
		"user":          userResponse(result.User),
	})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req refreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	result, apiErr := h.auth.Refresh(req.RefreshToken)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}

	response.OK(c, gin.H{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"token_type":    result.TokenType,
		"expires_in":    result.ExpiresIn,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	token = strings.TrimSpace(token)
	h.auth.Logout(token)
	response.NoContent(c)
}

func (h *Handler) Me(c *gin.Context) {
	user := middleware.UserFromContext(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	freshUser, ok := h.auth.Profile(user.ID)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	response.OK(c, userResponse(freshUser))
}

func userResponse(u *model.User) gin.H {
	return gin.H{
		"id":        u.ID,
		"role":      u.Role,
		"full_name": u.FullName,
		"email":     u.Email,
		"phone":     u.Phone,
		"status":    u.Status,
	}
}
