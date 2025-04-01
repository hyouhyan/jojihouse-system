package handler

import "jojihouse-entrance-system/internal/service"

type EntranceHandler struct {
	service *service.EntranceService
}

func NewEntranceHandler(service *service.EntranceService) *EntranceHandler {
	return &EntranceHandler{service: service}
}
