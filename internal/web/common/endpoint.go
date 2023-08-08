package common

import (
	"context"
	"net/http"

	apierrors "github.com/eurofurence/reg-room-service/internal/errors"
	"github.com/eurofurence/reg-room-service/internal/logging"
)

type (
	RequestHandler[Req any]  func(r *http.Request) (*Req, error)
	ResponseHandler[Res any] func(ctx context.Context, res *Res, w http.ResponseWriter) error
	Endpoint[Req, Res any]   func(ctx context.Context, request *Req) (*Res, error)
)

func CreateHandler[Req, Res any](endpoint Endpoint[Req, Res],
	requestHandler RequestHandler[Req],
	responseHandler ResponseHandler[Res],
) http.Handler {
	if endpoint == nil {
		panic("unable to set up service: no endpoint provided")
	}

	if requestHandler == nil {
		panic("unable to set up service: request handler must not be nil")
	}

	if responseHandler == nil {
		panic("unable to set up service: response handler must not be nil")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.LoggerFromContext(ctx)

		defer func() {
			err := r.Body.Close()
			if err != nil {
				logger.Error("Error when closing the request body. [error]: %v", err)
			}
		}()

		request, err := requestHandler(r)
		if err != nil {
			logger.Error("An error occurred while parsing the request. [error]: %v", err)

			if status := apierrors.AsAPIStatus(err); err != nil {
				SendHttpStatusErrorResponse(ctx, w, status)
				return
			}

			SendBadRequestResponse(ctx, w, "")
			return
		}

		response, err := endpoint(ctx, request)
		if err != nil {
			logger.Error("An error occurred during the request. [error]: %v", err)

			// In order to let the business logic decide, what kind of errors we want
			// to return, it makes sense to make use of a general service error type
			// which holds information about the error and the http status which should
			// be returned to the client.
			if status := apierrors.AsAPIStatus(err); status != nil {
				SendHttpStatusErrorResponse(ctx, w, status)
				return
			}

			// Fallback to internal server error if the original error could
			// not be determined properly.
			logger.Error("Service reported internal error: [error]: %v", err)
			SendInternalServerError(ctx, w, "")
			return
		}

		if err := responseHandler(ctx, response, w); err != nil {
			logger.Error("An error occurred during the handling of the response. [error]: %v", err)

			if status := apierrors.AsAPIStatus(err); err != nil {
				SendHttpStatusErrorResponse(ctx, w, status)
				return
			}

			SendInternalServerError(ctx, w, "")
		}
	})
}
