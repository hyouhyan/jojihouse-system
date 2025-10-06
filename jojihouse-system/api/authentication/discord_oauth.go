package authentication

import (
	"encoding/json"
	"fmt"
	"io"
	"jojihouse-system/api/model/response"
	"jojihouse-system/internal/service"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const DISCORD_API_BASEURL = "https://discordapp.com/api"

func Env_load() {
	err := godotenv.Load("./api/authentication/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

type DiscordAuthentication struct {
	userPortalService *service.UserPortalService
}

func NewDiscordAuthentication(userPortalService *service.UserPortalService) *DiscordAuthentication {
	return &DiscordAuthentication{userPortalService: userPortalService}
}

func (a *DiscordAuthentication) GetToken(code string) (token string, err error) {
	Env_load()

	// 色々定義
	values := url.Values{}
	values.Add("client_id", os.Getenv("CLIENT_ID"))
	values.Add("client_secret", os.Getenv("CLIENT_SECRET"))
	values.Add("grant_type", os.Getenv("GRANT_TYPE"))
	values.Add("redirect_uri", os.Getenv("REDIRECT_URL"))
	values.Add("code", code)

	// リクエストの作成
	req, err := http.NewRequest(
		"POST",
		DISCORD_API_BASEURL+"/oauth2/token",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("リクエストの作成に失敗しました: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Client作ってリクエスト投げる
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// ステータスコード確認
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code %d", res.StatusCode)
	}

	// Bodyの取得
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to get body %v", err)
	}

	// JSONをパース
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// tokenを抽出
	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token not found")
	}

	return token, nil
}

func (a *DiscordAuthentication) GetDiscordUserID(token string) (userID string, err error) {
	// Request組み立て
	req, err := http.NewRequest(
		"GET",
		DISCORD_API_BASEURL+"/users/@me",
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("リクエストの作成に失敗しました: %v", err)
	}

	// ヘッダー
	req.Header.Set("Authorization", "Bearer "+token)

	// Client作ってリクエスト投げる
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// ステータスコード確認
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code %d", res.StatusCode)
	}

	// Bodyの取得
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to get body %v", err)
	}

	// JSONをパース
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// idを抽出
	userID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("id not found")
	}

	_, err = strconv.Atoi(userID)
	if err != nil {
		return "", fmt.Errorf("discord id is not number: %v", err)
	}

	return userID, nil
}

func (a *DiscordAuthentication) GetHouseSystemUser(token string) (user *response.User, err error) {
	// Discord APIよりID取得
	discordUserID, err := a.GetDiscordUserID(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get discord user id: %v", err)
	}

	// HouseSystemからUser情報取得
	user, err = a.userPortalService.GetUserByDiscordID(discordUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by discord id: %v", err)
	}

	return user, nil
}
