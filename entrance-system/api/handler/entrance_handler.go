package handler

import (
	"jojihouse-entrance-system/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EntranceHandler struct {
	service *service.EntranceService
}

func NewEntranceHandler(service *service.EntranceService) *EntranceHandler {
	return &EntranceHandler{service: service}
}

type EntranceRequest struct {
	Barcode string `json:"barcode" binding:"required"`
	Type    string `json:"type" binding:"required,oneof=entry exit"`
}

// 入退室記録
func (h *EntranceHandler) RecordEntrance(c *gin.Context) {
	var req EntranceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Type == "entry" {
		err := h.service.EnterUser(req.Barcode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record entry"})
			return
		}
	} else {
		err := h.service.ExitUser(req.Barcode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record exit"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
