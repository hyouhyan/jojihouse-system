package handler

import (
	"jojihouse-entrance-system/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EntranceHandler struct {
	entranceService   *service.EntranceService
	userPortalService *service.UserPortalService
}

func NewEntranceHandler(entranceService *service.EntranceService, userPortalService *service.UserPortalService) *EntranceHandler {
	return &EntranceHandler{entranceService: entranceService, userPortalService: userPortalService}
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
		err := h.entranceService.EnterUser(req.Barcode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record entry"})
			return
		}
	} else {
		err := h.entranceService.ExitUser(req.Barcode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record exit"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// 在室ユーザー取得
func (h *EntranceHandler) GetCurrentUsers(c *gin.Context) {
	currentUsers, err := h.userPortalService.GetCurrentUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"current_users": currentUsers})

}
