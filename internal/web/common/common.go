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

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}

	return "ffffffff"
}

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

func GetGroups(ctx context.Context) []string {
	claims := GetClaims(ctx)
	if claims == nil || claims.Groups == nil {
		return []string{}
	}
	return claims.Groups
}

func HasGroup(ctx context.Context, group string) bool {
	for _, grp := range GetGroups(ctx) {
		if grp == group {
			return true
		}
	}
	return false
}

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
			aulogging.Logger.Ctx(ctx).Error().Printf("Could not encode response. [error]: %v", err)
		}
	}
}

func SendHTTPStatusErrorResponse(ctx context.Context, w http.ResponseWriter, status apierrors.APIStatus) {
	reqID := GetRequestID(ctx)
	w.WriteHeader(status.Status().Code)

	var detailValues url.Values
	details := status.Status().Details
	if details != "" {
		aulogging.Logger.Ctx(ctx).Debug().Printf("Request was not successful: [error]: %s", details)
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
		aulogging.Logger.Ctx(ctx).Debug().Printf("Request was not successful: [error]: %s", details)
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

func SendUnauthorizedResponse(ctx context.Context, w http.ResponseWriter, details string) {
	SendErrorWithStatusAndMessage(ctx, w, http.StatusUnauthorized, AuthUnauthorizedMessage, details)
}
