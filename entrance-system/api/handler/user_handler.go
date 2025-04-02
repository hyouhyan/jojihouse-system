package handler

import (
	"jojihouse-entrance-system/internal/request"
	"jojihouse-entrance-system/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userportalService      *service.UserPortalService
	adminManagementService *service.AdminManagementService
}

func NewUserHandler(userportalService *service.UserPortalService, adminManagementService *service.AdminManagementService) *UserHandler {
	return &UserHandler{userportalService: userportalService, adminManagementService: adminManagementService}
}

// ユーザー作成
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req request.CreateUserRequest
	// リクエストの解読
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := h.adminManagementService.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": res})
}

// ユーザー情報取得

// 全ユーザーの情報取得
