package handler

import (
	"jojihouse-entrance-system/api/model/request"
	"jojihouse-entrance-system/api/model/response"
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EntranceHandler struct {
	entranceService   *service.EntranceService
	userPortalService *service.UserPortalService
}

func NewEntranceHandler(entranceService *service.EntranceService, userPortalService *service.UserPortalService) *EntranceHandler {
	return &EntranceHandler{entranceService: entranceService, userPortalService: userPortalService}
}

// 入退室記録
func (h *EntranceHandler) RecordEntrance(c *gin.Context) {
	var req request.Entrance
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
	lastID := c.Query("last_id") // クエリパラメータから lastID を取得
	limitStr := c.Query("limit") // クエリパラメータから limit を取得

	// デフォルトの取得件数を設定（limit が指定されていなければ 10）
	limit := int64(10)
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = int64(parsedLimit)
		}
	}

	options := model.AccessLogFilter{
		Limit: limit,
	}

	accesLogs, err := h.userPortalService.GetAccessLogsByAnyFilter(lastID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_logs": accesLogs})
}

// アクセスログをユーザー指定で取得
func (h *EntranceHandler) GetAccessLogsByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	lastID := c.Query("last_id") // クエリパラメータから lastID を取得
	limitStr := c.Query("limit") // クエリパラメータから limit を取得

	// デフォルトの取得件数を設定（limit が指定されていなければ 10）
	limit := int64(10)
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = int64(parsedLimit)
		}
	}

	options := model.AccessLogFilter{
		Limit:  limit,
		UserID: userID,
	}

	accesLogs, err := h.userPortalService.GetAccessLogsByAnyFilter(lastID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_logs": accesLogs})
}
