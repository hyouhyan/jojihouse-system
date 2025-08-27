package handler

import (
	"jojihouse-system/api/authentication"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	discordAuthentication *authentication.DiscordAuthentication
	authentication        *authentication.TokenAuthentication
}

func NewAuthHandler(
	discordAuthentication *authentication.DiscordAuthentication,
	authentication *authentication.TokenAuthentication,
) *AuthHandler {
	return &AuthHandler{
		discordAuthentication: discordAuthentication,
		authentication:        authentication,
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
	user, err := h.discordAuthentication.GetHouseSystemUser(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get discord user id"})
		log.Print(err)
		return
	}

	// トークンの発行
	authToken, err := h.authentication.CreateJWTToken(*user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create authentication token"})
		log.Print(err)
		return
	}

	c.Header("Authorization", authToken)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
