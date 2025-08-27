package middleware

import (
	"jojihouse-system/api/authentication"
	"jojihouse-system/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	userPortalService   *service.UserPortalService
	tokenAuthentication *authentication.TokenAuthentication
}

func NewAuthMiddleware(userPortalService *service.UserPortalService, tokenAuthentication *authentication.TokenAuthentication) *AuthMiddleware {
	return &AuthMiddleware{
		userPortalService:   userPortalService,
		tokenAuthentication: tokenAuthentication,
	}
}

func (m AuthMiddleware) AuthSystemAdmin(c *gin.Context) {
	// Authorizationヘッダーからトークンを取得
	tokenString := c.GetHeader("Authorization")

	// Tokenの検証
	userID, err := m.tokenAuthentication.VerifyJWTToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		log.Printf("invalid token: %v", err)
		c.Abort()
		return
	}

	// ロール情報の検証
	isSysAdmin, err := m.userPortalService.IsSystemAdmin(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong user"})
		log.Printf("failed to check user is system admin: %v", err)
		c.Abort()
		return
	}

	if isSysAdmin {
		c.Next()
	} else {
		c.Abort()
	}
}
