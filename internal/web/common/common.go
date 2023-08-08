package common

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v4"

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

func SendBadRequestResponse(ctx context.Context, w http.ResponseWriter, details string) {
	SendResponseWithStatusAndMessage(
		w,
		http.StatusBadRequest,
		logging.GetRequestID(ctx),
		RequestParseErrorMessage,
		logging.LoggerFromContext(ctx),
		details)
}

func SendInternalServerError(ctx context.Context, w http.ResponseWriter, details string) {
	SendResponseWithStatusAndMessage(
		w,
		http.StatusInternalServerError,
		logging.GetRequestID(ctx),
		InternalErrorMessage,
		logging.LoggerFromContext(ctx),
		details)
}

func SendResponseWithStatusAndMessage(w http.ResponseWriter, status int, reqID string, message APIErrorMessage, logger logging.Logger, details string) {
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
