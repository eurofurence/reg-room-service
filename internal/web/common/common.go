package common

import (
	"context"
	"encoding/json"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"

	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
)

type (
	CtxKeyIDToken     struct{}
	CtxKeyAccessToken struct{}
	CtxKeyAPIKey      struct{}
	CtxKeyClaims      struct{}

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

func EncodeToJSON(ctx context.Context, w http.ResponseWriter, obj interface{}) {
	enc := json.NewEncoder(w)

	if obj != nil {
		err := enc.Encode(obj)
		if err != nil {
			aulogging.ErrorErrf(ctx, err, "Could not encode response. [error]: %v", err)
		}
	}
}

// SendHTTPStatusErrorResponse will send an api error
// which contains relevant information about the failed request to the client.
// The function will also set the http status according to the provided status.
func SendHTTPStatusErrorResponse(ctx context.Context, w http.ResponseWriter, status apierrors.APIStatus) {
	reqID := GetRequestID(ctx)
	w.WriteHeader(status.Status().Code)

	var detailValues url.Values
	details := status.Status().Details
	if details != "" {
		aulogging.Debugf(ctx, "Request was not successful: [error]: %s", details)
		detailValues = url.Values{"details": []string{details}}
	}

	apiErr := NewAPIError(reqID, status.Status().Message, detailValues)
	EncodeToJSON(ctx, w, apiErr)
}

// SendErrorWithStatusAndMessage will construct an api error
// which contains relevant information about the failed request to the client
// The function will also set the http status according to the provided status.
func SendErrorWithStatusAndMessage(ctx context.Context, w http.ResponseWriter, status int, message string, details string) {
	reqID := GetRequestID(ctx)
	w.WriteHeader(status)

	var detailValues url.Values
	if details != "" {
		aulogging.Debugf(ctx, "Request was not successful: [error]: %s", details)
		detailValues = url.Values{"details": []string{details}}
	}

	apiErr := NewAPIError(reqID, message, detailValues)
	EncodeToJSON(ctx, w, apiErr)
}

// EncodeWithStatus will attempt to encode the provided `value` into the
// response writer `w` and will write the status header.
// If the encoding fails, the http status will not be written to the response writer
// and the function will return an error instead.
func EncodeWithStatus[T any](status int, value *T, w http.ResponseWriter) error {
	err := json.NewEncoder(w).Encode(value)
	if err != nil {
		return errors.Wrap(err, "could not encode type into response buffer")
	}

	w.WriteHeader(status)

	return nil
}

// SendUnauthorizedResponse sends a standardized StatusUnauthorized response to the client.
func SendUnauthorizedResponse(ctx context.Context, w http.ResponseWriter, details string) {
	SendErrorWithStatusAndMessage(ctx, w, http.StatusUnauthorized, AuthUnauthorizedMessage, details)
}
