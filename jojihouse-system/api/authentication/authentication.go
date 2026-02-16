package authentication

import (
	"fmt"
	"jojihouse-system/internal/service"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenAuthentication struct {
	userPortalService *service.UserPortalService
}

func NewTokenAuthentication(userPortalService *service.UserPortalService) *TokenAuthentication {
	return &TokenAuthentication{userPortalService: userPortalService}
}

func (a *TokenAuthentication) CreateJWTToken(userID int) (tokenStr string, err error) {
	Env_load()

	// トークン発行
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(userID),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Stringに変換
	tokenStr, err = token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("failed to convert token to string: %v", err)
	}

	log.Print("[Token Authentication] Create JWT Token", tokenStr)

	return tokenStr, nil
}

func (a *TokenAuthentication) VerifyJWTToken(tokenStr string) (userID int, err error) {
	Env_load()

	// トークンの検証
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to verify jwt token: %v", err)
	} else if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Payloadの取得
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("failed to get claims")
	}

	// sub(userid)の抽出
	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, fmt.Errorf("failed to get User ID")
	}

	// intへ変換
	userID, err = strconv.Atoi(sub)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID: %w", err)
	}

	return userID, nil
}
