package handler

import (
	"jojihouse-entrance-system/api/model/request"
	"jojihouse-entrance-system/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userPortalService      *service.UserPortalService
	adminManagementService *service.AdminManagementService
}

func NewUserHandler(userPortalService *service.UserPortalService, adminManagementService *service.AdminManagementService) *UserHandler {
	return &UserHandler{userPortalService: userPortalService, adminManagementService: adminManagementService}
}

// ユーザー作成
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req request.CreateUser
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
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// URLパラメータから user_id を取得
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// サービス層でユーザー情報を取得
	user, err := h.userPortalService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// レスポンスを返す
	c.JSON(http.StatusOK, user)
}

// 全ユーザーの情報取得
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	res, err := h.userPortalService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
	}

	c.JSON(http.StatusOK, gin.H{"users": res})
}

// ユーザー情報を更新
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// URLパラメータから user_id を取得
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req request.UpdateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	h.adminManagementService.UpdateUser(userID, &req)

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// ユーザーを削除
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.adminManagementService.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// ロール取得
func (h *UserHandler) GetRolesByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	res, err := h.userPortalService.GetRolesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": res})
}
