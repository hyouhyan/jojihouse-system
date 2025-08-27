package handler

import (
	"jojihouse-system/api/authentication"
	"jojihouse-system/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userPortalService     *service.UserPortalService
	discordAuthentication *authentication.DiscordAuthentication
}

func NewAuthHandler(
	userPortalService *service.UserPortalService,
	discordAuthentication *authentication.DiscordAuthentication,
) *AuthHandler {
	return &AuthHandler{
		userPortalService:     userPortalService,
		discordAuthentication: discordAuthentication,
	}
}

func (h *AuthHandler) DiscordAuth(c *gin.Context) {
	// code query(パラメータ)の取得
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code query is necessary"})
	}

	// Tokenの取得
	token, err := h.discordAuthentication.GetToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get discord token"})
		log.Print(err)
	}

	log.Println(token)
}
