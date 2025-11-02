package handler

import (
	"errors"
	"jojihouse-system/api/model/request"
	"jojihouse-system/internal/model"
	"jojihouse-system/internal/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	adminManagementService *service.AdminManagementService
}

func NewPaymentHandler(
	adminManagementService *service.AdminManagementService,
) *PaymentHandler {
	return &PaymentHandler{
		adminManagementService: adminManagementService,
	}
}

func (h *PaymentHandler) GetAllPaymentLogs(c *gin.Context) {
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

	paymentLogs, err := h.adminManagementService.GetAllPaymentLogs(lastID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"title":  "Failed to get payment log",
			"detail": "支払いログの取得に失敗しました。",
		})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_logs": paymentLogs})
}

func (h *PaymentHandler) GetAllDeletedPaymentLogs(c *gin.Context) {
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

	paymentLogs, err := h.adminManagementService.GetAllDeletedPaymentLogs(lastID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"title":  "Failed to get deleted payment log",
			"detail": "削除済み支払いログの取得に失敗しました。",
		})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_logs": paymentLogs})
}

func (h *PaymentHandler) GetMonthlyPaymentLogs(c *gin.Context) {
	year := c.Query("year")
	month := c.Query("month")

	if year == "" || month == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid parameters",
			"detail": "yearとmonthが必要です。",
		})
		return
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid year format",
			"detail": "yearが誤っています。yearは整数で指定してください。",
		})
		return
	}
	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid month format",
			"detail": "monthが誤っています。monthは整数で、1~12の範囲で指定してください。",
		})
		return
	}

	paymentLogs, err := h.adminManagementService.GetMonthlyPaymentLogs(yearInt, monthInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"title":  "Failed to get monthly payment logs",
			"detail": "月間支払いログの取得に失敗しました。",
		})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"monthlyPaymentLogs": paymentLogs})
}

func (h *PaymentHandler) CreatePaymentLog(c *gin.Context) {
	var req request.Payment

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid request",
			"detail": "リクエストの形式が間違っています。",
		})
		log.Print(err)
		return
	}

	paymentLog := &model.PaymentLog{
		UserID:      req.UserID,
		Amount:      req.Amount,
		Description: req.Description,
		Payway:      req.Payway,
	}

	_, err := h.adminManagementService.CreatePaymentLog(paymentLog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"title":  "Failed to create payment log",
			"detail": "支払いの記録に失敗しました。",
		})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"time": paymentLog.Time.Format(time.RFC3339), "amount": paymentLog.Amount, "payway": paymentLog.Payway, "description": paymentLog.Description})
}

func (h *PaymentHandler) GetPaymentLogByID(c *gin.Context) {
	logID := c.Param("log_id")
	if logID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Log ID is required",
			"detail": "log_idパラメータが必要です。",
		})
		return
	}

	paymentLog, err := h.adminManagementService.GetPaymentLogByID(logID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPaymentLogNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"title":  "Payment log not found",
				"detail": "指定された支払いログが見つかりませんでした。",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "Failed to get payment log",
				"detail": "支払いログの取得に失敗しました。",
			})
		}
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_log": paymentLog})
}

func (h *PaymentHandler) DeletePaymentLog(c *gin.Context) {
	logID := c.Param("log_id")
	if logID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Log ID is required",
			"detail": "log_idパラメータが必要です。",
		})
		return
	}

	err := h.adminManagementService.DeletePaymentLog(logID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPaymentLogNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"title":  "Payment log not found",
				"detail": "指定された支払いログが見つかりませんでした。",
			})
		case errors.Is(err, model.ErrPaymentLogSeemsTicketPurchase):
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "Cannot delete Special log",
				"detail": "該当ログの削除には、特殊な処理が必要です。管理者に連絡してください。",
			})
		case errors.Is(err, model.ErrPaymentLogTooOldToDelete):
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "Payment log is too old to delete",
				"detail": "14日以上前のログは削除できません。どうしても削除したい場合は、管理者に連絡してください。",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "Failed to delete payment log",
				"detail": "支払いログの削除に失敗しました。",
			})
		}
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment log deleted successfully"})
}
