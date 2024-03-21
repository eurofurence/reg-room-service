package rbac

import (
	"context"

	"github.com/eurofurence/reg-room-service/internal/config"
	"github.com/eurofurence/reg-room-service/internal/web/common"
)

type CtxKeyValidator struct{}

type Validator struct {
	subject          string
	groups           []string
	isAdmin          bool
	isAPITokenCall   bool
	isRegisteredUser bool
}

func (i *Validator) IsAdmin() bool {
	return i.isAdmin
}

func (i *Validator) IsAPITokenCall() bool {
	return i.isAPITokenCall
}

func (i *Validator) IsRegisteredUser() bool {
	return i.isRegisteredUser && i.subject != ""
}

func (i *Validator) Subject() string {
	return i.subject
}

// ValidatorFromContext returns a Validator from the context.
// If the context doesn't contain a validator instance,
// a new instance will be returned instead.
func ValidatorFromContext(ctx context.Context) (*Validator, error) {
	validator, exists := ctx.Value(CtxKeyValidator{}).(*Validator)
	if exists {
		return validator, nil
	}

	return NewValidator(ctx)
}

// NewValidator returns a new instance of Validator.
// The function requires that the application config has been initialized.
func NewValidator(ctx context.Context) (*Validator, error) {
	conf, err := config.GetApplicationConfig()
	if err != nil {
		return nil, err
	}

	manager := &Validator{}
	if _, ok := ctx.Value(common.CtxKeyAPIKey{}).(string); ok {
		manager.isAPITokenCall = true
		return manager, nil
	}

	if claims, ok := ctx.Value(common.CtxKeyClaims{}).(*common.AllClaims); ok {
		manager.subject = claims.Subject
		manager.groups = claims.Groups

		manager.isRegisteredUser = true

		for _, group := range claims.Groups {
			if group == conf.Security.Oidc.AdminGroup && hasValidAdminHeader(ctx) {
				manager.isRegisteredUser = false
				manager.isAdmin = true
				break
			}
		}
	}

	return manager, nil
}

// TODO remove after 2FA is available
// See reference https://github.com/eurofurence/reg-payment-service/issues/57
func hasValidAdminHeader(ctx context.Context) bool {
	adminHeaderValue, ok := ctx.Value(common.CtxKeyAdminHeader{}).(string)
	if !ok {
		return false
	}

	// legacy system implementation requires check against constant value "available"
	return adminHeaderValue == "available"
}
