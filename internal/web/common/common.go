package common

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"

	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/logging"
)

type (
	CtxKeyIdToken     struct{}
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

func EncodeToJSON(w http.ResponseWriter, obj interface{}, logger logging.Logger) {
	enc := json.NewEncoder(w)

	if obj != nil {
		err := enc.Encode(obj)
		if err != nil {
			logger.Error("Could not encode response. [error]: %v", err)
		}
	}
}

func SendHttpStatusErrorResponse(ctx context.Context, w http.ResponseWriter, status apierrors.APIStatus) {
	logger := logging.LoggerFromContext(ctx)
	reqID := logging.GetRequestID(ctx)
	if reqID == "" {
		logger.Debug("request id is empty")
	}

	w.WriteHeader(status.Status().Code)

	var detailValues url.Values
	details := status.Status().Details
	if details != "" {
		logger.Debug("Request was not successful: [error]: %s", details)
		detailValues = url.Values{"details": []string{details}}
	}

	apiErr := NewAPIError(reqID, APIErrorMessage(status.Status().Message), detailValues)
	EncodeToJSON(w, apiErr, logger)
}

// SendErrorWithStatusAndMessage will construct an api error
// which contains relevant information about the failed request to the client
// The function will also set the http status according to the provided status
func SendErrorWithStatusAndMessage(w http.ResponseWriter, status int, reqID string, message APIErrorMessage, logger logging.Logger, details string) {
	if reqID == "" {
		logger.Debug("request id is empty")
	}

	w.WriteHeader(status)

	var detailValues url.Values
	if details != "" {
		logger.Debug("Request was not successful: [error]: %s", details)
		detailValues = url.Values{"details": []string{details}}
	}

	apiErr := NewAPIError(reqID, message, detailValues)
	EncodeToJSON(w, apiErr, logger)
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
