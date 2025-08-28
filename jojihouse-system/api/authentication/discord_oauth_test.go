package authentication

import (
	"fmt"
	"testing"
)

func TestGetToken_CodeからTokenを取得(t *testing.T) {
	discordAuthentication := NewDiscordAuthentication()

	token, err := discordAuthentication.GetToken("")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(token)
}

func TestGetUserID(t *testing.T) {
	discordAuthentication := NewDiscordAuthentication()

	id, err := discordAuthentication.GetUserID("")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(id)
}
