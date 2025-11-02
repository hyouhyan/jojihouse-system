package handler

import (
	"errors"
	"jojihouse-system/api/model/request"
	"jojihouse-system/api/model/response"
	"jojihouse-system/internal/model"
	"jojihouse-system/internal/service"
	"log"
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
// @Success 200 {object} response.User
// @Router /users [POST]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req request.CreateUser
	// リクエストの解読
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid request",
			"detail": "リクエストの形式が間違っています。",
		})
		log.Print(err)
		return
	}

	res, err := h.adminManagementService.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"title":  "Failed to create user",
			"detail": "ユーザの作成に失敗しました。",
		})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": res})
}

// @Summary ユーザー情報取得
// @Description 指定したユーザーの情報を取得します
// @Tags ユーザー管理
// @Produce json
// @Param user_id path int true "ユーザーID"
// @Success 200 {object} response.User
// @Router /users/{user_id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// URLパラメータから user_id を取得
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid user ID",
			"detail": "ユーザIDが誤っています。ユーザIDは整数で指定してください。",
		})
		log.Print(err)
		return
	}

	// サービス層でユーザー情報を取得
	user, err := h.userPortalService.GetUserByID(userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"title":  "User not found",
				"detail": "指定されたユーザが見つかりませんでした。",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "Failed to get user",
				"detail": "ユーザ情報の取得に失敗しました。",
			})
		}
		log.Print(err)
		return
	}

	// レスポンスを返す
	c.JSON(http.StatusOK, user)
}

// @Summary ユーザーの情報取得
// @Tags ユーザー管理
// @Produce json
// @Param barcode path int false "バーコード"
// @Success 200 {object} []response.User
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	// クエリパラメータからbarcodeを取得
	barcode := c.Query("barcode")

	if barcode != "" {
		// サービス層でユーザー情報を取得
		user, err := h.userPortalService.GetUserByBarcode(barcode)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrUserNotFound):
				c.JSON(http.StatusNotFound, gin.H{
					"title":  "User not found",
					"detail": "バーコードに対応するユーザが見つかりませんでした。",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"title":  "Failed to get user",
					"detail": "ユーザ情報の取得に失敗しました。",
				})
			}
			log.Print(err)
			return
		}

		// レスポンスを返す
		c.JSON(http.StatusOK, gin.H{"users": user})
		return
	}

	res, err := h.userPortalService.GetAllUsers()
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"title":  "User not found",
				"detail": "ユーザが見つかりませんでした。",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "Failed to get user",
				"detail": "ユーザ情報の取得に失敗しました。",
			})
		}
		log.Print(err)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid user ID",
			"detail": "ユーザIDが誤っています。ユーザIDは整数で指定してください。",
		})
		log.Print(err)
		return
	}

	var req request.UpdateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid request",
			"detail": "リクエストの形式が間違っています。",
		})
		log.Print(err)
		return
	}

	err = h.adminManagementService.UpdateUser(userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"title":  "User not found",
				"detail": "指定されたユーザが見つかりませんでした。",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "Failed to update user",
				"detail": "ユーザ情報の更新に失敗しました。",
			})
		}
		log.Print(err)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid user ID",
			"detail": "ユーザIDが誤っています。ユーザIDは整数で指定してください。",
		})
		log.Print(err)
		return
	}

	err = h.adminManagementService.DeleteUser(userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"title":  "User not found",
				"detail": "指定されたユーザが見つかりませんでした。",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "Failed to delete user",
				"detail": "ユーザの削除に失敗しました。",
			})
		}
		log.Print(err)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid user ID",
			"detail": "ユーザIDが誤っています。ユーザIDは整数で指定してください。",
		})
		log.Print(err)
		return
	}

	res, err := h.userPortalService.GetRolesByUserID(userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"title":  "User not found",
				"detail": "指定されたユーザが見つかりませんでした。",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "Failed to get user roles",
				"detail": "ユーザのロール情報の取得に失敗しました。",
			})
		}
		log.Print(err)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid user ID",
			"detail": "ユーザIDが誤っています。ユーザIDは整数で指定してください。",
		})
		log.Print(err)
		return
	}

	var req request.AddRole
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid request",
			"detail": "リクエストの形式が間違っています。",
		})
		log.Print(err)
		return
	}

	if err := h.adminManagementService.AddRoleToUser(userID, req.RoleID); err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"title":  "User not found",
				"detail": "指定されたユーザが見つかりませんでした。",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "Failed to add role to user",
				"detail": "ユーザへのロール追加に失敗しました。",
			})
		}
		log.Print(err)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid user ID",
			"detail": "ユーザIDが誤っています。ユーザIDは整数で指定してください。",
		})
		log.Print(err)
		return
	}

	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid role ID",
			"detail": "ロールIDが誤っています。ロールIDは整数で指定してください。",
		})
		log.Print(err)
		return
	}

	if err := h.adminManagementService.RemoveRoleFromUser(userID, roleID); err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"title":  "User not found",
				"detail": "指定されたユーザが見つかりませんでした。",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "Failed to remove role",
				"detail": "ユーザのロール削除に失敗しました。",
			})
		}
		log.Print(err)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid user ID",
			"detail": "ユーザIDが誤っています。ユーザIDは整数で指定してください。",
		})
		log.Print(err)
		return
	}

	lastID := c.Query("last_id") // クエリパラメータから lastID を取得

	remLogs, err := h.userPortalService.GetRemainingEntriesLogsByUserID(userID, lastID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"title":  "Failed to get remaining entries logs",
			"detail": "入場可能回数の変更ログの取得に失敗しました。",
		})
		log.Print(err)
		return
	}

	logs := response.Logs{
		RemainingEntriesLog: remLogs,
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func (h *UserHandler) ChangeRemainingEntries(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid user ID",
			"detail": "ユーザIDが誤っています。ユーザIDは整数で指定してください。",
		})
		log.Print(err)
		return
	}

	var req request.ChangeRemainingEntries
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Invalid request",
			"detail": "リクエストの形式が間違っています。",
		})
		log.Print(err)
		return
	}

	if req.Delta == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"title":  "Delta must be non-zero",
			"detail": "deltaは0以外の整数を指定してください。",
		})
		return
	}

	if req.Delta > 0 {
		_, err = h.adminManagementService.IncreaseRemainingEntries(userID, req.Delta, req.Readon, req.UpdatedBy)
	} else {
		_, err = h.adminManagementService.DecreaseRemainingEntries(userID, -req.Delta, req.Readon, req.UpdatedBy)
	}

	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"title":  "User not found",
				"detail": "指定されたユーザが見つかりませんでした。",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"title":  "Failed to change remaining entries",
				"detail": "入場可能回数の変更に失敗しました。",
			})
		}
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
