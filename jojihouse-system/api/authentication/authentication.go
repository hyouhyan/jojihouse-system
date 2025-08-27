package authentication

import "jojihouse-system/internal/service"

type Authentication struct {
	userPortalService *service.UserPortalService
}

func NewAuthentication(userPortalService *service.UserPortalService) *Authentication {
	return &Authentication{userPortalService: userPortalService}
}
