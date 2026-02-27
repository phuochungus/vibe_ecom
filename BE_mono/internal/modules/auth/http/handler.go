package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	authsvc "golf-store/be-mono/internal/modules/auth/service"
	"golf-store/be-mono/internal/platform/db"
	"golf-store/be-mono/internal/shared/middleware"
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

type userDTO struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

type loginResponseDTO struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	TokenType    string  `json:"token_type"`
	ExpiresIn    int     `json:"expires_in"`
	User         userDTO `json:"user"`
}

type refreshResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
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

	response.OK(c, loginResponseDTO{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
		ExpiresIn:    result.ExpiresIn,
		User:         userResponse(result.User),
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

	response.OK(c, refreshResponseDTO{
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

func userResponse(u *db.UserEntity) userDTO {
	if u == nil {
		return userDTO{}
	}
	return userDTO{
		ID:       u.ID,
		Role:     u.Role,
		FullName: u.FullName,
		Email:    u.Email,
		Phone:    u.Phone,
		Status:   u.Status,
	}
}
