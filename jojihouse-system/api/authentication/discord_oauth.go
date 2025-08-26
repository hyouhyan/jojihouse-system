package authentication

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const DISCORD_API_URL = "https://discordapp.com/api/oauth2/token"

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
		DISCORD_API_URL,
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
