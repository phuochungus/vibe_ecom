package http

import (
	usersvc "golf-store/be-mono/internal/modules/user/service"
)

type Handler struct {
	userSvc *usersvc.Service
}
