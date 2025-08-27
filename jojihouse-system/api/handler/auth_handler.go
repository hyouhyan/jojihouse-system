package handler

import (
	"jojihouse-system/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userPortalService *service.UserPortalService
}

func NewAuthHandler(
	userPortalService *service.UserPortalService,
) *AuthHandler {
	return &AuthHandler{
		userPortalService: userPortalService,
	}
}

func (h *AuthHandler) DiscordAuth(c *gin.Context) {
	code := c.Query("code")
}
