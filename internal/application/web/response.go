package web

import (
	"context"
	"encoding/json"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/eurofurence/reg-room-service/internal/application/common"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func EncodeToJSON(ctx context.Context, w http.ResponseWriter, obj interface{}) {
	enc := json.NewEncoder(w)

	if obj != nil {
		err := enc.Encode(obj)
		if err != nil {
			aulogging.ErrorErrf(ctx, err, "Could not encode response. [error]: %v", err)
		}
	}
}

// SendErrorResponse will send HTTPStatusErrorResponse if err is common.APIError.
//
// Otherwise sends internal server error.
func SendErrorResponse(ctx context.Context, w http.ResponseWriter, err error) {
	if err == nil {
		aulogging.ErrorErrf(ctx, err, "nil error in web layer")
		SendErrorWithStatusAndMessage(ctx, w, http.StatusInternalServerError, common.InternalErrorMessage, "an unspecified error occurred. Please check the logs - this is a bug")
		return
	}

	apiErr, ok := err.(common.APIError)
	if !ok {
		aulogging.ErrorErrf(ctx, err, "unwrapped error in web layer: %s", err.Error())
		SendErrorWithStatusAndMessage(ctx, w, http.StatusInternalServerError, common.InternalErrorMessage, "an unclassified error occurred. Please check the logs - this is a bug")
		return
	}
	SendAPIErrorResponse(ctx, w, apiErr)
}

// SendAPIErrorResponse will send an api error
// which contains relevant information about the failed request to the client.
// The function will also set the http status according to the provided status.
func SendAPIErrorResponse(ctx context.Context, w http.ResponseWriter, apiErr common.APIError) {
	aulogging.InfoErrf(ctx, apiErr, fmt.Sprintf("api response status %d: %v", apiErr.Status(), apiErr.Response()))
	for _, cause := range apiErr.InternalCauses() {
		aulogging.InfoErrf(ctx, cause, fmt.Sprintf("... caused by: %v", cause))
	}

	w.WriteHeader(apiErr.Status())

	EncodeToJSON(ctx, w, apiErr.Response())
}

// SendErrorWithStatusAndMessage will construct an api error
// which contains relevant information about the failed request to the client
// The function will also set the http status according to the provided status.
func SendErrorWithStatusAndMessage(ctx context.Context, w http.ResponseWriter, status int, message common.ErrorMessageCode, details string) {
	var detailValues url.Values
	if details != "" {
		aulogging.Debugf(ctx, "Request was not successful: [error]: %s", details)
		detailValues = url.Values{"details": []string{details}}
	}

	apiErr := common.NewAPIError(ctx, status, message, detailValues)
	SendAPIErrorResponse(ctx, w, apiErr)
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
	SendErrorWithStatusAndMessage(ctx, w, http.StatusUnauthorized, common.AuthUnauthorized, details)
}
