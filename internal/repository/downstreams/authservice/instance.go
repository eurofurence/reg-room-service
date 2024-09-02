package authservice

import (
	aulogging "github.com/StephanHCB/go-autumn-logging"

	"github.com/eurofurence/reg-room-service/internal/repository/config"
)

var activeInstance AuthService

func New(authServiceBaseUrl string, conf config.SecurityConfig) (AuthService, error) {
	if authServiceBaseUrl != "" {
		instance, err := newClient(authServiceBaseUrl, conf)
		activeInstance = instance
		return instance, err
	} else {
		aulogging.Logger.NoCtx().Warn().Printf("security.oidc.auth_service not configured. Will skip online userinfo checks (not useful for production!)")
		activeInstance = newMock()
		return activeInstance, nil
	}
}

func CreateMock() Mock {
	instance := newMock()
	activeInstance = instance
	return instance
}

func Get() AuthService {
	return activeInstance
}
