package handler

import (
	"jojihouse-entrance-system/api/model/request"
	"jojihouse-entrance-system/api/model/response"
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

// @Summary ユーザーを新規作成
// @Tags ユーザー管理
// @Produce json
// @Param user body request.CreateUser true "ユーザー情報"
// @Success 200 {object} response.UserResponse
// @Router /users/ [POST]
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

// @Summary ユーザー情報取得
// @Description 指定したユーザーの情報を取得します
// @Tags ユーザー管理
// @Produce json
// @Param user_id path int true "ユーザーID"
// @Success 200 {object} response.UserResponse
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

// @Summary 全ユーザーの情報取得
// @Tags ユーザー管理
// @Produce json
// @Success 200 {object} []response.UserResponse
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	res, err := h.userPortalService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
	}

	c.JSON(http.StatusOK, gin.H{"users": res})
}

// @Summary ユーザー情報を更新
// @Description 指定したユーザーの情報を更新します（部分更新）
// @Accept json
// @Tags ユーザー管理
// @Produce json
// @Param user_id path int true "ユーザーID"
// @Param user body request.UpdateUser true "更新するユーザー情報（部分的に送信可能）"
// @Success 200 {object} map[string]string
// @Router /users/{user_id} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
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

	err = h.adminManagementService.UpdateUser(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// @Summary ユーザーを削除
// @Description 指定したユーザーを削除します
// @Tags ユーザー管理
// @Produce json
// @Param user_id path int true "ユーザーID"
// @Success 200 {object} map[string]string
// @Router /users/{user_id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.adminManagementService.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// @Summary 指定ユーザーのロールを取得
// @Description 指定したユーザーが持つロールを取得します
// @Tags ユーザー管理
// @Produce json
// @Param user_id path int true "ユーザーID"
// @Success 200 {object} []response.Role
// @Router /users/{user_id}/roles [get]
func (h *UserHandler) GetRolesByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	res, err := h.userPortalService.GetRolesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": res})
}

// @Summary 指定ユーザーにロールを追加
// @Description 指定したユーザーにロールを追加します
// @Accept json
// @Tags ユーザー管理
// @Produce json
// @Param user_id path int true "ユーザーID"
// @Param role body request.AddRole true "追加するロールのID"
// @Success 200 {object} map[string]string
// @Router /users/{user_id}/roles [post]
func (h *UserHandler) AddRoleToUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req request.AddRole
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.adminManagementService.AddRoleToUser(userID, req.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// @Summary 指定ユーザーのロールを削除
// @Description 指定したユーザーからロールを削除します
// @Tags ユーザー管理
// @Produce json
// @Param user_id path int true "ユーザーID"
// @Param role_id path int true "ロールID"
// @Success 200 {object} map[string]string
// @Router /users/{user_id}/roles/{role_id} [delete]
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// @Summary ユーザーのログを取得
// @Description 指定したユーザーの各ログを取得します
// @Tags ユーザー管理
// @Produce json
// @Param user_id path int true "ユーザーID"
// @Success 200 {object} response.Logs
// @Router /users/{user_id}/logs [get]
func (h *UserHandler) GetUserLogs(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	lastID := c.Query("last_id") // クエリパラメータから lastID を取得

	remLogs, err := h.userPortalService.GetRemainingEntriesLogsByUserID(userID, lastID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get remaining entries logs"})
		return
	}

	logs := response.Logs{
		RemainingEntriesLog: remLogs,
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
