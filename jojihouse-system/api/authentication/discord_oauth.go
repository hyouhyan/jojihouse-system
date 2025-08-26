package authentication

import "jojihouse-system/internal/service"

type authEnv struct {
	clientID     string
	clientSecret string
	grantType    string
	refirectUsr  string
	code         string
}

type DiscordAuthentication struct {
	userPortalService *service.UserPortalService
}

func NewDiscordAuthentication(userPortalService *service.UserPortalService) *DiscordAuthentication {
	return &DiscordAuthentication{userPortalService: userPortalService}
}
