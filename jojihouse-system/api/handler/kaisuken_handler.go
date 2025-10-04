package handler

import (
	"jojihouse-system/api/model/request"
	"jojihouse-system/internal/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type KaisukenHandler struct {
	userPortalService      *service.UserPortalService
	adminManagementService *service.AdminManagementService
}

func NewKaisukenHandler(
	userPortalService *service.UserPortalService,
	adminManagementService *service.AdminManagementService,
) *KaisukenHandler {
	return &KaisukenHandler{
		userPortalService:      userPortalService,
		adminManagementService: adminManagementService,
	}
}

// @Summary 回数券購入
// @Tags 回数券
// @Description ユーザーが回数券を購入します
// @Accept json
// @Produce json
// @Param kaisuken body request.BuyKaisuken true "回数券購入のリクエストデータ"
// @Success 200 {object} map[string]string
// @Router /kaisuken [post]
func (h *KaisukenHandler) BuyKaisuken(c *gin.Context) {
	var req request.BuyKaisuken
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		log.Print(err)
		return
	}

	paymentLog, err := h.adminManagementService.BuyKaisuken(req.UserID, req.Receiver, req.Amount, req.Count, req.Payway, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record kaisuken purchase"})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"time": paymentLog.Time.Format(time.RFC3339), "amount": paymentLog.Amount, "payway": paymentLog.Payway, "description": paymentLog.Description})
}
