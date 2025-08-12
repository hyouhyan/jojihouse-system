package handler

import (
	"errors"
	"jojihouse-entrance-system/api/model/request"
	"jojihouse-entrance-system/api/model/response"
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EntranceHandler struct {
	entranceService   *service.EntranceService
	userPortalService *service.UserPortalService
}

func NewEntranceHandler(entranceService *service.EntranceService, userPortalService *service.UserPortalService) *EntranceHandler {
	return &EntranceHandler{entranceService: entranceService, userPortalService: userPortalService}
}

// @Summary 入退室記録
// @Tags エントランス(入退室)管理
// @Description ユーザーの入室または退室を記録します
// @Accept json
// @Produce json
// @Param entrance body request.Entrance true "入退室データ"
// @Success 200 {object} response.Entrance
// @Router /entrance [post]
func (h *EntranceHandler) RecordEntrance(c *gin.Context) {
	var req request.Entrance
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		log.Print(err)
		return
	}

	var response response.Entrance
	var err error

	// autoの場合
	if req.Type == "auto" {
		// ユーザーの在室を確認
		isCurrent, err := h.userPortalService.IsCurrentUserByBarcode(req.Barcode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check current user"})
			log.Print(err)
			return
		}
		if isCurrent {
			req.Type = "exit"
		} else {
			req.Type = "entry"
		}
	}

	switch req.Type {
	case "entry":
		response, err = h.entranceService.EnterUser(req.Barcode)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrUserNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record entry"})
			}
			log.Print(err)
			return
		}
	case "exit":
		response, err = h.entranceService.ExitUser(req.Barcode)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrUserNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record exit"})
			}
			log.Print(err)
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type. Must be 'entry', 'exit', or 'auto'"})
		log.Print("Invalid type. " + req.Type)
		return
	}

	c.JSON(http.StatusOK, gin.H{"entrance_log": response})
}

// @Summary 在室ユーザー取得
// @Tags エントランス(入退室)管理
// @Description 現在ハウス内にいるユーザーの一覧を取得します
// @Produce json
// @Success 200 {object} []response.User
// @Router /entrance/current [get]
func (h *EntranceHandler) GetCurrentUsers(c *gin.Context) {
	currentUsers, err := h.userPortalService.GetCurrentUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current users"})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"current_users": currentUsers})
}

// @Summary アクセスログを取得
// @Tags エントランス(入退室)管理
// @Description すべてのユーザーの入退室ログを取得します
// @Produce json
// @Param last_id query string false "前回のログID（ページネーション用）"
// @Param limit query int false "取得するログの件数（デフォルト10）"
// @Param date query string false "対象日（YYYY-MM-DD形式）"
// @Success 200 {object} []response.AccessLog
// @Router /entrance/logs [get]
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

	dateStr := c.Query("date")

	var date time.Time
	var err error

	// 日付のバリデーション（YYYY-MM-DD の形式）
	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			log.Print(err)
			return
		}
	}

	options := model.AccessLogFilter{
		Limit:     limit,
		DayBefore: date,
		DayAfter:  date,
	}

	accessLogs, err := h.userPortalService.GetAccessLogsByAnyFilter(lastID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access log"})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_logs": accessLogs})
}

// @Summary アクセスログをユーザー指定で取得
// @Tags エントランス(入退室)管理
// @Description 指定したユーザーの入退室ログを取得します
// @Produce json
// @Param user_id path int true "ユーザーID"
// @Param last_id query string false "前回のログID（ページネーション用）"
// @Param limit query int false "取得するログの件数（デフォルト10）"
// @Success 200 {object} []response.AccessLog
// @Router /entrance/logs/{user_id} [get]
func (h *EntranceHandler) GetAccessLogsByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		log.Print(err)
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

	accessLogs, err := h.userPortalService.GetAccessLogsByAnyFilter(lastID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access log"})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_logs": accessLogs})
}

// @Summary 入退室ログを手動で追加
// @Tags エントランス(入退室)管理
// @Description 手動で入退室ログを追加します
// @Accept json
// @Produce json
// @Param entrance body request.FixedAccessLog true "入退室データ"
// @Success 200 {object} response.Entrance
// @Router /entrance/fix [post]
func (h *EntranceHandler) CreateFixedAccessLog(c *gin.Context) {
	var req request.FixedAccessLog
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		log.Print(err)
		return
	}

	// バーコード、UserID、会員番号のいずれかが埋まってるか確認
	if req.Barcode == nil && req.UserID == nil && req.Number == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either barcode, user_id or number must be provided"})
		log.Print("Either barcode, user_id or number must be provided")
		return
	}

	err := h.entranceService.CreateFixedAccessLog(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create fixed access log"})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Fixed access log created successfully"})
}
