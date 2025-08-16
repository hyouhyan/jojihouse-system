package handler

import (
	"jojihouse-system/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	service *service.UserPortalService
}

func NewRoleHandler(service *service.UserPortalService) *RoleHandler {
	return &RoleHandler{service: service}
}

// @Summary 全ロール取得
// @Tags ロール管理
// @Description 存在する全ロールをまとめて取得
// @Produce json
// @Success 200 {object} []response.Role
// @Router /roles [get]
func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	res, err := h.service.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get roles"})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": res})
}
