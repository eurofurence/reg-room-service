package common

import (
	"context"
	"net/http"

	"github.com/eurofurence/reg-room-service/internal/apierrors"
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
		reqID := logging.GetRequestID(ctx)
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
			SendBadRequestResponse(w, reqID, logger, "")
			return
		}

		response, err := endpoint(ctx, request)
		if err != nil {
			logger.Error("An error occurred during the request. [error]: %v", err)

			// check if the error is a `StatusError`
			if status := apierrors.AsAPIStatus(err); status != nil {

				// TODO enhance
				switch {
				case apierrors.IsBadRequestError(err):
					SendBadRequestResponse(w, reqID, logger, status.Status().Details)
				case apierrors.IsUnauthorizedError(err):
					SendUnauthorizedResponse(w, reqID, logger, status.Status().Details)
				case apierrors.IsForbiddenError(err):
					SendForbiddenResponse(w, reqID, logger, status.Status().Details)
				case apierrors.IsNotFoundError(err):
					SendStatusNotFoundResponse(w, reqID, logger, status.Status().Details)
				case apierrors.IsConflictError(err):
					SendConflictResponse(w, reqID, logger, status.Status().Details)
				case apierrors.IsInternalServerError(err):
					SendInternalServerError(w,
						reqID,
						logger,
						status.Status().Details,
					)
				}

				return
			}

			// do not propagate internal errors to the client.
			// check the logs for errors - and use metrics later on
			logger.Error("Service reported internal error: [error]: %v", err)
			SendInternalServerError(w, reqID, logger, "")
			return
		}

		if err := responseHandler(ctx, response, w); err != nil {
			logger.Error("An error occurred during the handling of the response. [error]: %v", err)
			SendInternalServerError(w, reqID, logger, "")
			return
		}
	})
}
