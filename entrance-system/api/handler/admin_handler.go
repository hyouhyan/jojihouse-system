package handler

import "jojihouse-entrance-system/internal/service"

type AdminHandler struct {
	service *service.AdminManagementService
}

func NewAdminHandler(service *service.AdminManagementService) *AdminHandler {
	return &AdminHandler{service: service}
}
