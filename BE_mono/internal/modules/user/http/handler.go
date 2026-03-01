package http

import (
	"net/http"

	usersvc "golf-store/be-mono/internal/modules/user/service"
	"golf-store/be-mono/internal/shared/response"

	"github.com/gin-gonic/gin"

	middleware "golf-store/be-mono/internal/shared/middleware"
)

type Handler struct {
	userSvc *usersvc.Service
}

func New(userSvc *usersvc.Service) *Handler {
	return &Handler{userSvc: userSvc}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.GET("/users/me", h.GetCurrentUser)
	rg.GET("/users/:id", h.GetUserByID)
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
	user := c.MustGet(middleware.ContextUserKey)
	response.OK(c, user)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	user, err := h.userSvc.GetUserByID(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found", nil)
		return
	}

	response.OK(c, user)
}
