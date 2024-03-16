package common

import (
	"context"
	"net/http"

	"github.com/eurofurence/reg-room-service/internal/logging"
)

type CtxKeyRequestURL struct{}

type (
	RequestHandler[Req any]  func(r *http.Request, w http.ResponseWriter) (*Req, error)
	ResponseHandler[Res any] func(ctx context.Context, res *Res, w http.ResponseWriter) error
	Endpoint[Req, Res any]   func(ctx context.Context, request *Req, w http.ResponseWriter) (*Res, error)
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

		ctx = context.WithValue(ctx, CtxKeyRequestURL{}, r.URL)

		defer func() {
			err := r.Body.Close()
			if err != nil {
				logger.Error("Error when closing the request body. [error]: %v", err)
			}
		}()

		request, err := requestHandler(r, w)
		if err != nil {
			logger.Error("An error occurred while parsing the request. [error]: %v", err)
		}

		response, err := endpoint(ctx, request, w)
		if err != nil {
			logger.Error("An error occurred during the request. [error]: %v", err)
			return
		}

		if err := responseHandler(ctx, response, w); err != nil {
			logger.Error("An error occurred during the handling of the response. [error]: %v", err)
		}
	})
}
