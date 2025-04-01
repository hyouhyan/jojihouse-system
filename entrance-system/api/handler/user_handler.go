package handler

import "jojihouse-entrance-system/internal/service"

type UserHandler struct {
	service *service.UserPortalService
}

func NewUserHandler(service *service.UserPortalService) *UserHandler {
	return &UserHandler{service: service}
}
