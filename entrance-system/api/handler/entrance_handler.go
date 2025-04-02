package handler

import (
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/response"
	"jojihouse-entrance-system/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	var response response.EntranceResponse
	var err error

	if req.Type == "entry" {
		response, err = h.entranceService.EnterUser(req.Barcode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record entry"})
			return
		}
	} else {
		response, err = h.entranceService.ExitUser(req.Barcode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record exit"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"entrance_log": response})
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

// アクセスログを取得
func (h *EntranceHandler) GetAccessLogs(c *gin.Context) {
	lastID := primitive.NilObjectID
	options := model.AccessLogFilter{
		Limit: 10,
	}

	accesLogs, err := h.userPortalService.GetAccessLogsByAnyFilter(lastID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest access log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_logs": accesLogs})
}
