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

// GetUserByID ユーザー情報取得
// @Summary ユーザー情報取得
// @Description 指定したユーザーの情報を取得します
// @Tags users
// @Produce json
// @Param user_id path int true "ユーザーID"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{user_id} [get]
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

// ロール追加
func (h *UserHandler) AddRoleToUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req request.AddRole
	// リクエストの解読
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ロール追加してエラーハンドリング
	if err := h.adminManagementService.AddRoleToUser(userID, req.RoleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not add role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// ロール削除
func (h *UserHandler) RemoveRoleFromUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	if err := h.adminManagementService.RemoveRoleFromUser(userID, roleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not remove role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
