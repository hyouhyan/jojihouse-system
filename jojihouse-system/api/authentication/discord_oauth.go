package authentication

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const DISCORD_API_BASEURL = "https://discordapp.com/api"

type DiscordAuthentication struct {
	// userPortalService *service.UserPortalService
}

// func NewDiscordAuthentication(userPortalService *service.UserPortalService) *DiscordAuthentication {
func NewDiscordAuthentication() *DiscordAuthentication {
	// return &DiscordAuthentication{userPortalService: userPortalService}
	return &DiscordAuthentication{}
}

func (a *DiscordAuthentication) GetToken(code string) (token string, err error) {
	// 色々定義
	values := url.Values{}
	values.Add("client_id", CLIENT_ID)
	values.Add("client_secret", CLIENT_SECRET)
	values.Add("grant_type", GRANT_TYPE)
	values.Add("redirect_uri", REDIRECT_URL)
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

func (a *DiscordAuthentication) GetUserID(token string) (userID int, err error) {
	// Request組み立て
	req, err := http.NewRequest(
		"GET",
		DISCORD_API_BASEURL+"/users/@me",
		nil,
	)
	if err != nil {
		return 0, fmt.Errorf("リクエストの作成に失敗しました: %v", err)
	}

	// ヘッダー
	req.Header.Set("Authorization", "Bearer "+token)

	// Client作ってリクエスト投げる
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	// ステータスコード確認
	if res.StatusCode != 200 {
		return 0, fmt.Errorf("status code %d", res.StatusCode)
	}

	// Bodyの取得
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to get body %v", err)
	}

	// JSONをパース
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to parse response: %v", err)
	}

	// idを抽出
	userIDStr, ok := result["id"].(string)
	if !ok {
		return 0, fmt.Errorf("id not found")
	}

	userID, err = strconv.Atoi(userIDStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse id: %v", err)
	}

	return userID, nil
}
