package handler

import (
	"fmt"
	"jojihouse-system/api/authentication"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	discordAuthentication *authentication.DiscordAuthentication
}

func NewAuthHandler(
	discordAuthentication *authentication.DiscordAuthentication,
) *AuthHandler {
	return &AuthHandler{
		discordAuthentication: discordAuthentication,
	}
}

func (h *AuthHandler) DiscordAuth(c *gin.Context) {
	// code query(パラメータ)の取得
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code query is necessary"})
		return
	}

	// Tokenの取得
	token, err := h.discordAuthentication.GetToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get discord token"})
		log.Print(err)
		return
	}

	// UserIDの取得
	discordUserID, err := h.discordAuthentication.GetUserID(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get discord user id"})
		log.Print(err)
		return
	}

	fmt.Println(discordUserID)

	c.JSON(http.StatusOK, gin.H{"message": "success", "discord_user_id": discordUserID})
}
