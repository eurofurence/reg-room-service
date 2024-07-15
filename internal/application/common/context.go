package common

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
)

type (
	CtxKeyIDToken     struct{}
	CtxKeyAccessToken struct{}
	CtxKeyAPIKey      struct{}
	CtxKeyClaims      struct{}

	CtxKeyRequestID  struct{}
	CtxKeyRequestURL struct{}

	// TODO Remove after legacy system was replaced with 2FA
	// See reference https://github.com/eurofurence/reg-payment-service/issues/57
	CtxKeyAdminHeader struct{}
)

type CustomClaims struct {
	EMail         string   `json:"email"`
	EMailVerified bool     `json:"email_verified"`
	Groups        []string `json:"groups,omitempty"`
	Name          string   `json:"name"`
}

type AllClaims struct {
	jwt.RegisteredClaims
	CustomClaims
}

const RequestIDKey = "RequestIDKey"

// GetRequestID extracts the request ID from the context.
//
// It originally comes from a header with the request, or is rolled while processing
// the request.
func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}

	return "ffffffff"
}

// GetClaims extracts all jwt token claims from the context.
func GetClaims(ctx context.Context) *AllClaims {
	claims := ctx.Value(CtxKeyClaims{})
	if claims == nil {
		return nil
	}

	allClaims, ok := claims.(*AllClaims)
	if !ok {
		return nil
	}

	return allClaims
}

// GetGroups extracts the groups from the jwt token that came with the request
// or from the groups retrieved from userinfo, if using authorization token.
//
// In either case the list is filtered by relevant groups (if reg-auth-service is configured).
func GetGroups(ctx context.Context) []string {
	claims := GetClaims(ctx)
	if claims == nil || claims.Groups == nil {
		return []string{}
	}
	return claims.Groups
}

// HasGroup checks that the user has a group.
func HasGroup(ctx context.Context, group string) bool {
	for _, grp := range GetGroups(ctx) {
		if grp == group {
			return true
		}
	}
	return false
}

// GetSubject extracts the subject field from the jwt token or the userinfo response, if using
// an authorization token.
func GetSubject(ctx context.Context) string {
	claims := GetClaims(ctx)
	if claims == nil {
		return ""
	}
	return claims.Subject
}
