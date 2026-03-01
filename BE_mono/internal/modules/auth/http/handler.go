package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"golf-store/be-mono/internal/modules/auth/dto"
	"golf-store/be-mono/internal/platform/db"
	"golf-store/be-mono/internal/shared/middleware"
	"golf-store/be-mono/internal/shared/response"

	authsvc "golf-store/be-mono/internal/modules/auth/service"
)

type Handler struct {
	authsvc *authsvc.Service
}

func New(auth *authsvc.Service) *Handler {
	return &Handler{authsvc: auth}
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
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	result, apiErr := h.authsvc.Login(req.Identifier, req.Password)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}

	response.OK(c, dto.LoginResponseDTO{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
		ExpiresIn:    result.ExpiresIn,
		User:         userResponse(result.User),
	})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	result, apiErr := h.authsvc.Refresh(req.RefreshToken)
	if apiErr != nil {
		response.Error(c, apiErr.Status, apiErr.Code, apiErr.Message, apiErr.Details)
		return
	}

	response.OK(c, dto.RefreshResponseDTO{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
		ExpiresIn:    result.ExpiresIn,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	token = strings.TrimSpace(token)
	h.authsvc.Logout(token)
	response.NoContent(c)
}

func (h *Handler) Me(c *gin.Context) {
	user := middleware.UserFromContext(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	freshUser, ok := h.authsvc.Profile(user.ID)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	response.OK(c, userResponse(freshUser))
}

func userResponse(u *db.UserEntity) dto.UserDTO {
	if u == nil {
		return dto.UserDTO{}
	}
	return dto.UserDTO{
		ID:       u.ID,
		Role:     u.Role,
		FullName: u.FullName,
		Email:    u.Email,
		Phone:    u.Phone,
		Status:   u.Status,
	}
}
