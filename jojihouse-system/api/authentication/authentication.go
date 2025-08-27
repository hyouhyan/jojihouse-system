package authentication

import (
	"fmt"
	"jojihouse-system/internal/service"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Authentication struct {
	userPortalService *service.UserPortalService
}

func NewAuthentication(userPortalService *service.UserPortalService) *Authentication {
	return &Authentication{userPortalService: userPortalService}
}

func (a *Authentication) CreateJWTToken(userID int) (tokenStr string, err error) {
	Env_load()

	// トークン発行
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Stringに変換
	tokenStr, err = token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("failed to convert token to string: %v", err)
	}

	return tokenStr, nil
}

func (a *Authentication) VerifyJWTToken(tokenStr string) (ok bool, err error) {
	Env_load()

	// トークンの検証
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return false, fmt.Errorf("failed to verify jwt token: %v", err)
	} else if !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}
